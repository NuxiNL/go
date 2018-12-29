package runtime

// +build cloudabi linux
// +build 386 amd64 arm arm64 ppc64 ppc64le

const (
	_PT_LOAD    = 1 /* Loadable program segment */
	_PT_DYNAMIC = 2 /* Dynamic linking information */

	_DT_NULL     = 0          /* Marks end of dynamic section */
	_DT_HASH     = 4          /* Dynamic symbol hash table */
	_DT_STRTAB   = 5          /* Address of string table */
	_DT_SYMTAB   = 6          /* Address of symbol table */
	_DT_GNU_HASH = 0x6ffffef5 /* GNU-style dynamic symbol hash table */
	_DT_VERSYM   = 0x6ffffff0
	_DT_VERDEF   = 0x6ffffffc

	_VER_FLG_BASE = 0x1 /* Version definition of file itself */

	_SHN_UNDEF = 0 /* Undefined section */

	_SHT_DYNSYM = 11 /* Dynamic linker symbol table */

	_STT_FUNC = 2 /* Symbol is a code object */

	_STT_NOTYPE = 0 /* Symbol type is not specified */

	_STB_GLOBAL = 1 /* Global symbol */
	_STB_WEAK   = 2 /* Weak symbol */

	_EI_NIDENT = 16
)
