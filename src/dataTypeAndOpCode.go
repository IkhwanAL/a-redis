package src

/*
@see {https://rdb.fnordig.de/file_format.html}
*/

// Byte Representing As Data Type in RDB File Format
var STRING_DATATYPE_FLAG = '\x00'

// OpCode
var START_METADATA = '\xFA'       // Auxiliary fields
var START_DB_SECTION = '\xFE'     // Database Selector
var START_HASHTABEL_INFO = '\xFB' // ResizeDB
var EOF = '\xFF'
var EXPIRETIMEMS = '\xFC'
