package testutils


const BUBBLESORT_SRC_2 = `
/* sort the array */
func bubbleSort(arr: integer[], size: integer) -> void
{
  let n: integer;
  let i: integer;
  let j: integer;
  let temp: integer;
  n = size;
  i = 0;
  j = 0;
  temp = 0;
  while (i < n-1) {
    while (j < n-i-1) {
      if (arr[j] > arr[j+1])
        then {
          // swap temp and arr[i]
          temp = arr[j];
          arr[j] = arr[j+1];
          arr[j+1] = temp;
        } else ;
        j = j+1;
      };
    i = i+1;
  };
}

/* print the array */
// func printArray(arr: integer[], size: integer) -> void   // Changed function name to match function call in main
func printarray(arr: integer[], size: integer) -> void
{
  let n: integer;
  let i: integer;
  n = size;
  i = 0;
  while (i<n) {
    write(arr[i]);
      i = i+1;
  };
}

// main funtion to test above
func main() -> void
{
  let arr: integer[7];
  arr[0] = 64;
  arr[1] = 34;
  arr[2] = 25;
  arr[3] = 12;
  arr[4] = 22;
  arr[5] = 11;
  arr[6] = 90;
  printarray(arr, 7);
  bubbleSort(arr, 7);
  printarray(arr, 7);
}
`

// INLINED FILE: `polynomial.src` with some modifications related to type
// checking
const POLYNOMIAL_SRC_2 = `
// ====== struct declarations ====== //
struct POLYNOMIAL {
	public func evaluate(x: float) -> float;
};

struct LINEAR inherits POLYNOMIAL {
	private let a: float;
	private let b: float;
	public  func build(A: float, B: float) -> LINEAR;
	public  func evaluate(x: float) -> float;
};

struct QUADRATIC inherits POLYNOMIAL {
	private let a: float;
	private let b: float;
	private let c: float;
	public  func build(A: float, B: float, C: float) -> QUADRATIC;
	public  func evaluate(x: float) -> float;
};

// ====== struct implementations ====== //
impl POLYNOMIAL {
  func evaluate(x: float) -> float
  {
    // return (0);
    return (0.0);
  }
}

impl QUADRATIC {
  func evaluate(x: float) -> float
  {
    let result: float;
    //Using Horner's method
    result = a;
    result = result * x + b;
    result = result * x + c;
    return (result);
  }

  func build(A: float, B: float, C: float) -> QUADRATIC
  {
    let new_function: QUADRATIC ;
    new_function.a = A;
    new_function.b = B;
    new_function.c = C;
    return (new_function);
  }
}

impl LINEAR {
  func build(A: float, B: float) -> LINEAR
  {
    let new_function: LINEAR;
    new_function.a = A;
    new_function.b = B;
    return (new_function);
  }
  func evaluate(x: float) -> float
  {
    let result: float;
    result = 0.0;
    result = a * x + b;
    return (result);
  }
}

// ====== main ====== //
func main() -> void
{
  let f1: LINEAR;
  let f2: QUADRATIC;
  // let counter: integer;  // This is invalid, should be float
  let counter: float;
  // f1 = f1.build(2, 3.5);  // This is invalid because 2 should be a float: 2.0
  f1 = f1.build(2.0, 3.5);
  f2 = f2.build(-2.0, 1.0, 0.0);
  counter = 1.0;

  while(counter <= 10.0)
  {
    write(counter);
    write(f1.evaluate(counter));
    write(f2.evaluate(counter));
  };
}
`


// INLINED FILE: `polynomial-with-errors-2.src`
const POLYNOMIAL_WITH_ERRORS_2_SRC = `
// ====== struct declarations ====== //

// Removing the ID after the struct keyword
struct {
	public func evaluate(x: float) -> float;
};

// Removing the 'inherits'
struct LINEAR  POLYNOMIAL {

    // Removing visibility modifier
	let a: float;

	private let b: float;
	public  func build(A: float, B: float) -> LINEAR;

    // Removing visibility modifier
	func evaluate(x: float) -> float;
};

struct QUADRATIC inherits POLYNOMIAL {
	private let a: float;
	private let b: float;
	private let c: float;
	public  func build(A: float, B: float, C: float) -> QUADRATIC;
	public  func evaluate(x: float) -> float;
};

// ====== struct implementations ====== //
impl POLYNOMIAL {
  func evaluate(x: float) -> float
  {
    return (0);
  }
}

impl QUADRATIC {

  // Removing both the arrow '->' and the return type 'float'
  func evaluate(x: float)
  {
    let result: float;
    //Using Horner's method
    result = a;
    result = result * x + b;
    result = result * x + c;

    // Removing the 'return (result);' - this is not a syntax error according to the grammar
  }
  func build(A: float, B: float, C: float) -> QUADRATIC
  {
    let new_function: QUADRATIC ;
    new_function.a = A;
    new_function.b = B;
    new_function.c = C;
    return (new_function);
  }
}

impl LINEAR {
  func build(A: float, B: float) -> LINEAR
  {
    let new_function: LINEAR;
    new_function.a = A;
    new_function.b = B;
    return (new_function);
  }
  func evaluate(x: float) -> float
  {
    let result: float;
    result = 0.0;
    result = a * x + b;
    return (result);
  }
}

// ====== main ====== //
func main() -> void
{
  let f1: LINEAR;
  let f2: QUADRATIC;
  let counter: integer;
  f1 = f1.build(2, 3.5);
  f2 = f2.build(-2.0, 1.0, 0.0);
  counter = 1;

  while(counter <= 10)
  {
    write(counter);
    write(f1.evaluate(counter));
    write(f2.evaluate(counter));
  };
}
`

// INLINED FILE: `polynomial-with-errors.src`
const POLYNOMIAL_WITH_ERRORS_SRC = `
// ====== struct declarations ====== //
struct POLYNOMIAL {
    // Missing visibility modifier
	func evaluate(x: float) -> float;
};

struct LINEAR inherits POLYNOMIAL {
    // Missing id after let
	private let : float;
	private let b: float;

    // Missing function return type
	public  func build(A: float, B: float) -> ;
	public  func evaluate(x: float) -> float;
};

struct QUADRATIC inherits POLYNOMIAL {
	private let a: float;

    // Missing type
	private let b;

	private let c: float;

    // Replaced arrow with colon
	public  func build(A: float, B: float, C: float): QUADRATIC;

	public  func evaluate(x: float) -> float;
};

// ====== struct implementations ====== //
impl POLYNOMIAL {
  func evaluate(x: float) -> float
  {
    return (0);
  }
}

impl QUADRATIC {
  func evaluate(x: float) -> float
  {
    let result: float;
    //Using Horner's method
    result = a;
    result = result * x + b;
    result = result * x + c;
    return (result);
  }
  func build(A: float, B: float, C: float) -> QUADRATIC
  {
    let new_function: QUADRATIC ;
    new_function.a = A;
    new_function.b = B;
    new_function.c = C;
    return (new_function);
  }
}

impl LINEAR {
  func build(A: float, B: float) -> LINEAR
  {
    let new_function: LINEAR;
    new_function.a = A;
    new_function.b = B;
    return (new_function);
  }
  func evaluate(x: float) -> float
  {
    let result: float;
    result = 0.0;
    result = a * x + b;
    return (result);
  }
}

// ====== main ====== //
func main() -> void
{
  let f1: LINEAR;
  let f2: QUADRATIC;
  let counter: integer;
  f1 = f1.build(2, 3.5);
  f2 = f2.build(-2.0, 1.0, 0.0);
  counter = 1;

  while(counter <= 10)
  {
    write(counter);
    write(f1.evaluate(counter));
    write(f2.evaluate(counter));
  };
}
`

// INLINED FILE: `polynomial.src`
const POLYNOMIAL_SRC = `
// ====== struct declarations ====== //
struct POLYNOMIAL {
	public func evaluate(x: float) -> float;
};

struct LINEAR inherits POLYNOMIAL {
	private let a: float;
	private let b: float;
	public  func build(A: float, B: float) -> LINEAR;
	public  func evaluate(x: float) -> float;
};

struct QUADRATIC inherits POLYNOMIAL {
	private let a: float;
	private let b: float;
	private let c: float;
	public  func build(A: float, B: float, C: float) -> QUADRATIC;
	public  func evaluate(x: float) -> float;
};

// ====== struct implementations ====== //
impl POLYNOMIAL {
  func evaluate(x: float) -> float
  {
    return (0);
  }
}

impl QUADRATIC {
  func evaluate(x: float) -> float
  {
    let result: float;
    //Using Horner's method
    result = a;
    result = result * x + b;
    result = result * x + c;
    return (result);
  }
  func build(A: float, B: float, C: float) -> QUADRATIC
  {
    let new_function: QUADRATIC ;
    new_function.a = A;
    new_function.b = B;
    new_function.c = C;
    return (new_function);
  }
}

impl LINEAR {
  func build(A: float, B: float) -> LINEAR
  {
    let new_function: LINEAR;
    new_function.a = A;
    new_function.b = B;
    return (new_function);
  }
  func evaluate(x: float) -> float
  {
    let result: float;
    result = 0.0;
    result = a * x + b;
    return (result);
  }
}

// ====== main ====== //
func main() -> void
{
  let f1: LINEAR;
  let f2: QUADRATIC;
  let counter: integer;
  f1 = f1.build(2, 3.5);
  f2 = f2.build(-2.0, 1.0, 0.0);
  counter = 1;

  while(counter <= 10)
  {
    write(counter);
    write(f1.evaluate(counter));
    write(f2.evaluate(counter));
  };
}
`

// INLINED FILE: `bubblesort.src`
const BUBBLESORT_SRC = `
/* sort the array */
func bubbleSort(arr: integer[], size: integer) -> void
{
  let n: integer;
  let i: integer;
  let j: integer;
  let temp: integer;
  n = size;
  i = 0;
  j = 0;
  temp = 0;
  while (i < n-1) {
    while (j < n-i-1) {
      if (arr[j] > arr[j+1])
        then {
          // swap temp and arr[i]
          temp = arr[j];
          arr[j] = arr[j+1];
          arr[j+1] = temp;
        } else ;
        j = j+1;
      };
    i = i+1;
  };
}

/* print the array */
func printArray(arr: integer[], size: integer) -> void
{
  let n: integer;
  let i: integer;
  n = size;
  i = 0;
  while (i<n) {
    write(arr[i]);
      i = i+1;
  };
}

// main funtion to test above
func main() -> void
{
  let arr: integer[7];
  arr[0] = 64;
  arr[1] = 34;
  arr[2] = 25;
  arr[3] = 12;
  arr[4] = 22;
  arr[5] = 11;
  arr[6] = 90;
  printarray(arr, 7);
  bubbleSort(arr, 7);
  printarray(arr, 7);
}
`

// INLINED FILE: `unterminatedcomments.src`
const UNTERMINATEDCOMMENTS_SRC = `
// this is an inline comment

/* this is a single line block comment
`

// INLINED FILE: `unterminatedcomments.src`
const UNTERMINATEDCOMMENTS2_SRC = `
/* this is an imbricated
/* block comment
`

// INLINED FILE: `lexpositivegrading.src`
const LEX_POSITIVE_GRADING_SRC = `
==	+	|	(	;	if 	public	read
<>	-	&	)	,	then	private	write
<	*	!	{	.	else	func	return
>	/		}	:	integer	var	self
<=	=		[	::	float	struct	inherits
>=			]	->	void	while	let
						func	impl





0
1
10
12
123
12345

1.23
12.34
120.34e10
12345.6789e-123

abc
abc1
a1bc
abc_1abc
abc1_abc

// this is an inline comment

/* this is a single line block comment */

/* this is a
multiple line
block comment
*/

/* this is an imbricated
/* block comment
*/
*/




`

// INLINED FILE: `lexnegativegrading.src`
const LEX_NEGATIVE_GRADING_SRC = `
@ # $ ' \ ~

00
01
010
0120
01230
0123450

01.23
012.34
12.340
012.340

012.34e10
12.34e010

_abc
1abc
_1abc

`

// INLINED FILE: `helloworld.src`
const LEX_HELLOWORLD_SRC = `
/*
This is an imaginary program with a made up syntax.

Let us see how the parser handles it...
*/

// C-style struct
struct Student {
    float age;
    integer id;
};

public func main() {

    // x is my integer variable
    let x = 10;

    /*
    y is equal to x
    */
    var y = x;

    // Equality check
    if (y == x) then {
        var out integer[10] = {x, y, 69, 200, 89};
        write(out);  // Assume we have a function called 'write'
    }
}
`

// INLINED FILE: `strings.src`
const LEX_STRINGS_SRC = `
var x = "this is not valid";
`

// INLINED FILE: `somethingelse.src`
const LEX_SOMETHING_ELSE_SRC = `
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("https://www.google.com")
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	headers, err := json.MarshalIndent(resp.Header, "", "    ")
	if err != nil {
		log.Fatalf("Failed to serialize response headers: %v", err)
	}
	fmt.Println(string(headers))

	bod, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read body")
	}
	defer resp.Body.Close()

	fmt.Println(string(bod))
}
`
