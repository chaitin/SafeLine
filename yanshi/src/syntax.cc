#include "syntax.hh"

void stmt_free(Stmt* stmt)
{
  while (stmt) {
    auto x = stmt->next;
    delete stmt;
    stmt = x;
  }
}
