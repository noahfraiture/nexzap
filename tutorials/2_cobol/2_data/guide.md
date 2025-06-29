# Data Division

The Data Division is where you set up your variables and reserve memory for your COBOL program. It's like laying out the storage bins before you start cooking up logic in the Procedure Division. Here's a quick example:

```cobol
DATA DIVISION.
WORKING-STORAGE SECTION.
01 WS-NAME    PIC X(25).
01 WS-CLASS   PIC 9(2)  VALUE '10'.
```

## Sections

The Data Division is split into sections, with the **Working-Storage Section** being the star of the show for defining variables used in your program. Other sections, like the **File Section** (for file-related data) or **Linkage Section** (for passing data between programs), pop up depending on your needs. Each section keeps things organized, so your data knows exactly where to chill.

## Data Types

COBOL's data types are defined using the **Picture Clause (PIC)**, which spells out what kind of data a variable can hold:
- **Numeric (9)**: For numbers, like `PIC 9(3)` for a 3-digit number (e.g., 123).
- **Alphabetic (A)**: For letters A-Z and spaces, like `PIC A(6)` for a 6-letter name.
- **Alphanumeric (X)**: For any mix of digits, letters, or special characters, like `PIC X(10)` for a 10-character ID.
- **Special Symbols**: Use `V` for an implied decimal (e.g., `PIC 9(2)V9(2)` for 12.34), `S` for signed numbers, or `P` for assumed decimals in calculations.

You can also initialize variables with the **Value Clause**, like `PIC 9(2) VALUE 42`, to give them a starting value. Data can be **elementary** (single variables) or **group** items (collections of variables under one name), organized with **Level Numbers** (01 for groups, 02-49 for elementary items). For extra flair, **Redefines** lets you reuse the same memory for different data layouts, and **Usage Clause** tweaks how data is stored (e.g., `COMP-3` for packed decimals).

It's a bit old-school, but this setup gives you tight control over your data—perfect for those massive business datasets COBOL loves to crunch!
