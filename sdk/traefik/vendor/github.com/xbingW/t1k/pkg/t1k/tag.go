package t1k

type Tag byte

const (
	TAG_HEADER       Tag = 0x01
	TAG_BODY         Tag = 0x02
	TAG_EXTRA        Tag = 0x03
	TAG_RSP_HEADER   Tag = 0x11
	TAG_RSP_BODY     Tag = 0x12
	TAG_RSP_EXTRA    Tag = 0x13
	TAG_VERSION      Tag = 0x20
	TAG_ALOG         Tag = 0x21
	TAG_STAT         Tag = 0x22
	TAG_EXTRA_HEADER Tag = 0x23
	TAG_EXTRA_BODY   Tag = 0x24
	TAG_CONTEXT      Tag = 0x25
	TAG_COOKIE       Tag = 0x26
	TAG_WEB_LOG      Tag = 0x27
	TAG_USER_DATA    Tag = 0x28
	TAG_BOT_QUERY    Tag = 0x29
	TAG_DELAY        Tag = 0x2b
	TAG_FORWARD      Tag = 0x2c
	TAG_BOT_BODY     Tag = 0x30

	MASK_TAG   Tag = 0x3f
	MASK_FIRST Tag = 0x40
	MASK_LAST  Tag = 0x80
)

func (t Tag) Last() bool {
	return t&MASK_LAST != 0
}

func (t Tag) First() bool {
	return t&MASK_FIRST != 0
}

func (t Tag) Strip() Tag {
	return t & MASK_TAG
}

func (t Tag) Byte() byte {
	return byte(t)
}
