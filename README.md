# bitcoin-inscription-parser

bitcoin-inscription-parser is a tool which helps to parse bitcoin inscriptions 
from transactions. Any inscription content wrapped in `OP_FALSE OP_IF … OP_ENDIF`
using data pushes can be correctly parsed.

The tool supports single or multiple inscriptions in all input of the transaction.
