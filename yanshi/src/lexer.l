%{
#include "lexer_helper.hh"
#include "option.hh"
#include "parser.hh"
#include "syntax.hh"

#include <limits.h>
#include <stdlib.h>
#include <string>
#include <string.h>
#include <unicode/utf8.h>
using namespace std;

#define YY_USER_ACTION                      \
  do {                                      \
    yylloc->start = yyget_extra(yyscanner); \
    yylloc->end = yylloc->start + yyleng;   \
    yyset_extra(yylloc->end, yyscanner);    \
  } while (0);

static string tmp_bracket, tmp_str;
static long tmp_str_pos;
static bool semicolon;

static long invalid_escape(YYSTYPE* yylval, const char* text)
{
  yylval->errmsg = aprintf("invalid \\-escape: %s", text);
  return INVALID_CHARACTER;
}

static int invalid_escape_octonary(YYSTYPE* yylval, const char* text)
{
  yylval->errmsg = aprintf("invalid number after \\-escape: %s", text);
  return INVALID_CHARACTER;
}

static int invalid_escape_x(YYSTYPE* yylval, const char* text)
{
  yylval->errmsg = aprintf("invalid number after \\x-escape: %s", text);
  return INVALID_CHARACTER;
}

static int invalid_escape_u(YYSTYPE* yylval, const char* text)
{
  yylval->errmsg = aprintf("invalid number after \\u-escape: %s", text);
  return INVALID_CHARACTER;
}

static int invalid_escape_U(YYSTYPE* yylval, const char* text)
{
  yylval->errmsg = aprintf("invalid number after \\U-escape: %s", text);
  return INVALID_CHARACTER;
}

static void unexpected_eof(YYSTYPE* yylval, const char* token_end)
{
  yylval->errmsg = aprintf("missing %s at end of file", token_end);
}

static void unexpected_newline(YYSTYPE* yylval, const char* token_end)
{
  yylval->errmsg = aprintf("missing %s at end of line", token_end);
}

static int unexpected_codepoint(YYSTYPE* yylval)
{
  yylval->errmsg = aprintf("cannot use Unicode codepoints");
  return INVALID_CHARACTER;
}

extern "C" int raw_yywrap(yyscan_t yyscanner)
{
  semicolon = false;
  return 1;
}
%}

%option yywrap noinput
%option reentrant
%option extra-type="long"
%option bison-bridge bison-locations
%option prefix="raw_yy"
%option stack

%x EXPECT_CODE
%x AFTER_ACTION_OP
%x AFTER_EXPORT
%x IN_BRACE
%x IN_CODE
%x IN_COMMENT
%x IN_BRACKET
%x IN_BRACKET_FIRST
%x IN_LINE_COMMENT
%s IN_PAREN
%x IN_Q_STRING
%x IN_QQ_STRING

D			[0-9]
H			[0-9A-Fa-f]
L			[a-zA-Z_\x80-\xff]

%%

"::" return COLONCOLON;
".." return DOTDOT;
"&&" return AMPERAMPER;
";" if (semicolon) return '\n';
[-~!&*=+,.?|{}:] return yytext[0];
"action" yy_push_state(EXPECT_CODE, yyscanner); return ACTION;
"as" return AS;
"c++" yy_push_state(EXPECT_CODE, yyscanner); return CPP;
"epsilon" return EPSILON;
"export" yy_push_state(AFTER_EXPORT, yyscanner); return EXPORT;
"import" return IMPORT;
"intact" return INTACT;
"semicolon" semicolon = true;
"nosemicolon" semicolon = false;
{L}({L}|{D})* yylval->str = new string(yytext); return IDENT;
{D}+ yylval->integer = atol(yytext); return INTEGER;
"#define" return PREPROCESS_DEFINE;

"#" yy_push_state(IN_LINE_COMMENT, yyscanner);
"//" yy_push_state(IN_LINE_COMMENT, yyscanner);
<IN_LINE_COMMENT>{
  "\n" yy_pop_state(yyscanner); unput('\n'); yyset_extra(yylloc->end-1, yyscanner);
  <<EOF>> yy_pop_state(yyscanner);
  . {}
}

"/*" yy_push_state(IN_COMMENT, yyscanner);
<IN_COMMENT>{
  "*/" yy_pop_state(yyscanner);
  <<EOF>> yy_pop_state(yyscanner);
  .|\n {}
}

"(" yy_push_state(IN_PAREN, yyscanner); return '(';
")" {
  if (YY_START != IN_PAREN) {
    unexpected_newline(yylval, ")");
    return INVALID_CHARACTER;
  }
  yy_pop_state(yyscanner);
  return ')';
}

"[" yy_push_state(IN_BRACKET_FIRST, yyscanner); return '[';
<IN_BRACKET_FIRST>{
  "^" BEGIN IN_BRACKET; return '^';
  [^-\\\]\n] { yy_pop_state(yyscanner); yy_push_state(IN_BRACKET, yyscanner); yylval->integer = yytext[0]; return CHAR; }
  "-" { yylval->integer = '-'; return CHAR; }
}
<IN_BRACKET>{
  "]" {
    yy_pop_state(yyscanner);
    if (YY_START == INITIAL || YY_START == IN_PAREN)
      return ']';
  }
  [^-\\\]\n] yylval->integer = yytext[0]; return CHAR;
}
<IN_BRACKET_FIRST,IN_BRACKET>{
  \\[0-7]+ {
    BEGIN IN_BRACKET;
    long c = strtol(yytext+1, NULL, 8);
    if (UCHAR_MAX < c)
      return invalid_escape_octonary(yylval, yytext);
    yylval->integer = c;
    return CHAR;
  }
  \\u[0-9a-fA-F]+ {
    BEGIN IN_BRACKET;
    if (opt_bytes)
      return unexpected_codepoint(yylval);
    long c = strtol(yytext+2, NULL, 16), len = 0;
    if (UINT16_MAX < c)
      return invalid_escape_u(yylval, yytext);
    yylval->integer = c;
    return CHAR;
  }
  \\U[0-9a-fA-F]+ {
    BEGIN IN_BRACKET;
    if (opt_bytes)
      return unexpected_codepoint(yylval);
    long c = strtol(yytext+2, NULL, 16), len = 0;
    if (MAX_CODEPOINT < c)
      return invalid_escape_U(yylval, yytext);
    yylval->integer = c;
    return CHAR;
  }
  \\x[0-9a-fA-F]+ {
    BEGIN IN_BRACKET;
    long c = strtol(yytext+2, NULL, 16);
    if (UCHAR_MAX < c)
      return invalid_escape_x(yylval, yytext);
    yylval->integer = c;
    return CHAR;
  }
  \\a BEGIN IN_BRACKET; yylval->integer = '\a'; return CHAR;
  \\b BEGIN IN_BRACKET; yylval->integer = '\b'; return CHAR;
  \\f BEGIN IN_BRACKET; yylval->integer = '\f'; return CHAR;
  \\n BEGIN IN_BRACKET; yylval->integer = '\n'; return CHAR;
  \\r BEGIN IN_BRACKET; yylval->integer = '\r'; return CHAR;
  \\t BEGIN IN_BRACKET; yylval->integer = '\t'; return CHAR;
  \\v BEGIN IN_BRACKET; yylval->integer = '\v'; return CHAR;
  \\. BEGIN IN_BRACKET; yylval->integer = yytext[1]; return CHAR;
  - BEGIN IN_BRACKET; return '-';
  "\n" unexpected_newline(yylval, "]"); return INVALID_CHARACTER;
  <<EOF>> yy_pop_state(yyscanner); unexpected_eof(yylval, "]");
}

<AFTER_EXPORT>{ // optional 'BRACED_CODE' to specify extra parameters
  "intact" yy_pop_state(yyscanner); return INTACT;
  {L}({L}|{D})* yy_pop_state(yyscanner); yylval->str = new string(yytext); return IDENT;
  "{" BEGIN IN_CODE; tmp_bracket.clear();
  [ \t\n] {}
  <<EOF>> yy_pop_state(yyscanner);
}

<EXPECT_CODE>{
  {L}({L}|{D})* yylval->str = new string(yytext); return IDENT;
  "{" BEGIN IN_CODE; tmp_bracket.clear();
  [ \t\n] {}
  <<EOF>> yy_pop_state(yyscanner);
}

[>@%$] yy_push_state(AFTER_ACTION_OP, yyscanner); return yytext[0];
<AFTER_ACTION_OP>{
  -?{D}+ yylval->integer = atol(yytext); return INTEGER;
  {L}({L}|{D})* yy_pop_state(yyscanner); yylval->str = new string(yytext); return IDENT;
  "{" BEGIN IN_CODE; tmp_bracket.clear();
  [ \t\n]+ {}
  <<EOF>> yy_pop_state(yyscanner);
  . yylval->errmsg = strdup("invalid character"); return INVALID_CHARACTER;
}
<IN_CODE>{
  "'" { tmp_bracket += '\''; yy_push_state(IN_Q_STRING, yyscanner); }
  "\"" { tmp_bracket += '"'; yy_push_state(IN_QQ_STRING, yyscanner); }
  "{" { tmp_bracket += '{'; yy_push_state(IN_CODE, yyscanner); }
  "}" {
    yy_pop_state(yyscanner);
    if (YY_START == INITIAL || YY_START == IN_PAREN) {
      yylval->str = new string(tmp_bracket);
      return BRACED_CODE;
    } else
      tmp_bracket += '}';
  }
  .|"\n" tmp_bracket += yytext[0];
  <<EOF>> yy_pop_state(yyscanner); unexpected_eof(yylval, "}");
}

' tmp_str.clear(); tmp_str_pos = yylloc->start; yy_push_state(IN_Q_STRING, yyscanner);
"\"" tmp_str.clear(); tmp_str_pos = yylloc->start; yy_push_state(IN_QQ_STRING, yyscanner);
<IN_Q_STRING>{
  ' {
    yy_pop_state(yyscanner);
    if (YY_START == INITIAL || YY_START == IN_PAREN) {
      yylval->str = new string(tmp_str);
      yylloc->start = tmp_str_pos;
      return STRING_LITERAL;
    }
    tmp_bracket += yytext;
  }
  <<EOF>> yy_pop_state(yyscanner); unexpected_eof(yylval, "'");
}
<IN_QQ_STRING>{
  "\"" {
    yy_pop_state(yyscanner);
    if (YY_START == INITIAL || YY_START == IN_PAREN) {
      yylval->str = new string(tmp_str);
      yylloc->start = tmp_str_pos;
      return STRING_LITERAL;
    }
    tmp_bracket += yytext;
  }
  <<EOF>> yy_pop_state(yyscanner); unexpected_eof(yylval, "\"");
}

<IN_Q_STRING,IN_QQ_STRING>{
  \\[0-7]+ {
    long c = strtol(yytext+1, NULL, 8);
    if (UCHAR_MAX < c)
      return invalid_escape_octonary(yylval, yytext);
    tmp_str.push_back(c);
    tmp_bracket += yytext;
  }
  \\x[0-9a-fA-F]+ {
    long c = strtol(yytext+2, NULL, 16);
    if (UCHAR_MAX < c)
      return invalid_escape_x(yylval, yytext);
    tmp_str.push_back(c);
    tmp_bracket += yytext;
  }
  \\u[0-9a-fA-F]+ {
    char s[4];
    long c = strtol(yytext+2, NULL, 16), len = 0;
    if (UINT16_MAX < c)
      return invalid_escape_u(yylval, yytext);
    U8_APPEND_UNSAFE(s, len, c);
    tmp_str.insert(tmp_str.end(), s, s+len);
    tmp_bracket += yytext;
  }
  \\U[0-9a-fA-F]+ {
    char s[4];
    long c = strtol(yytext+2, NULL, 16), len = 0;
    if (MAX_CODEPOINT < c)
      return invalid_escape_U(yylval, yytext);
    U8_APPEND_UNSAFE(s, len, c);
    tmp_str.insert(tmp_str.end(), s, s+len);
    tmp_bracket += yytext;
  }
  \\a tmp_str += '\a'; tmp_bracket += yytext;
  \\b tmp_str += '\b'; tmp_bracket += yytext;
  \\f tmp_str += '\f'; tmp_bracket += yytext;
  \\n tmp_str += '\n'; tmp_bracket += yytext;
  \\r tmp_str += '\r'; tmp_bracket += yytext;
  \\t tmp_str += '\t'; tmp_bracket += yytext;
  \\v tmp_str += '\v'; tmp_bracket += yytext;
  \\[\n\"\'?\\] tmp_str += yytext[1]; tmp_bracket += yytext;
  \\. return invalid_escape(yylval, yytext);
  .|\n tmp_str += yytext[0]; tmp_bracket += yytext[0];
}

\\\n {}
"\n" if (YY_START == INITIAL && ! semicolon) return '\n';
[ \t]+ {}
. { yylval->errmsg = strdup("invalid character"); return INVALID_CHARACTER; }
