# What is Cobol?

Common Business Oriented Language (COBOL), is a very old language used in legacy system of finance mostly. It was designed for easy readability so that even non-technical people can read it and understand it (in theory, in practice the code can be really messy), and its capacity to easily manage file and high amount of data. Here is a small example with its english-like syntax :
```cobol
IDENTIFICATION DIVISION.
PROGRAM-ID. HELLO.

PROCEDURE DIVISION.
   DISPLAY 'Hello World'.
STOP RUN.
```
# Structure

Cobol is a programming language very different from modern programming language. A program is divided in 4 divisions : *Identification*, *Environment*, *Data* and *Procedure*.
1. **Identification Division**: Specifies the program's metadata, such as its name, author, and purpose, providing basic information about the program.
2. **Environment Division**: Defines the hardware and software environment, including computer configuration and input/output devices used by the program.
3. **Data Division**: Declares all data items, variables, and file structures, outlining how data is organized and stored within the program.
4. **Procedure Division**: Contains the executable code, including logic and instructions, that defines the program's operations and processing steps.

## Pros
- Can handle huge volume of data with advanced file handling capabilities
- High safety and precision design for business case
- English-like syntax making it explicit and easily readable

## Cons
- Limited eco-system: no one uses COBOL anymore except in specific businesses. The environment isn't expanding, and there are few libraries and frameworks
- Verbose syntax: COBOL's lengthy, rigid syntax, with extensive keywords and strict formatting rules, makes code harder to read and write compared to concise modern languages
- It's weird, the logic is weird. Totally subjective and biased, but it feels heavy and difficult for no good reason.

## Why Learn Cobol?

Cobol isn't used anymore to build new stuffs because people like Cobol. Cobol is used in the bank system because it's too complicated and costly to migrate to a more modern programming language. It doesn't mean Cobol doesn't have strength, it does. But except for curiosity, learning Cobol will not bring you professional value except if you go in financial system.
