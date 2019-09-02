interface LInterface {
    Name: string;
    readonly x: number;
    Sex?: string;
}

interface Square extends LInterface {
    sideLength: number;
}



let ro: ReadonlyArray<number> = [1, 2, 3, 4];

class Clock implements LInterface {
    Name: string = "666"
    x: number = 666
    constructor(n: string, m: number) {
    }

    static standardGreeting = "Hello, there";
    greet() {
        return Clock.standardGreeting;
    }
}

interface Counter {
    (start: number): string;
    interval: number;
    reset(): void;
}

// Hybrid Types #
function getCounter(): Counter {
    let counter = (function (start: number) { }) as Counter;
    counter.interval = 123;
    counter.reset = function () { };
    return counter;
}

let c = getCounter();
c(10);
c.reset()
c.interval = 5.0;

let someValue: any = "this is a string";
let strLength: number = (<string>someValue).length;

let myAdd: (x: number, y: number) => number =
    function (x: number, y: number): number { return x + y; };

function buildName(firstName: string, lastName?: string) {
    if (lastName)
        return firstName + " " + lastName;
    else
        return firstName;
}

let result1 = buildName("Bob");

function buildName2(firstName: string, lastName = "Smith") {
    return firstName + " " + lastName;
}

// Rest Parameters 
function buildName3(firstName: string, ...restOfName: string[]) {
    return firstName + " " + restOfName.join(" ");
}

// employeeName will be "Joseph Samuel Lucas MacKinzie"
let employeeName = buildName3("Joseph", "Samuel", "Lucas", "MacKinzie");

let deck = {
    suits: ["hearts", "spades", "clubs", "diamonds"],
    cards: Array(52),    
    createCardPicker: function() {
        return function() {
            let pickedCard = Math.floor(Math.random() * 52);
            let pickedSuit = Math.floor(pickedCard / 13);

            return {suit: this.suits[pickedSuit], card: pickedCard % 13};
        }
    }
}

let cardPicker = deck.createCardPicker();
let pickedCard = cardPicker();

alert("card: " + pickedCard.card + " of " + pickedCard.suit);

function loggingIdentity<T>(arg: T[]): T[] {
    console.log(arg.length);  // Array has a .length, so no more error
    return arg;
}

interface GenericIdentityFn {
    <T>(arg: T): T;
}

enum FileAccess {
    // constant members
    None,
    Read    = 1 << 1,
    Write   = 1 << 2,
    ReadWrite  = Read | Write,
    // computed member
    G = "123".length
}

let list = [4, 5, 6];

for (let i in list) {
    console.log(i); // "0", "1", "2",
}

for (let i of list) {
    console.log(i); // "4", "5", "6"
}

// import { StringValidator } from "./StringValidator";
// export * from "./StringValidator"; // exports 'StringValidator' interface
// export * from "./ZipCodeValidator";  // exports 'ZipCodeValidator' and const 'numberRegexp' class
// export * from "./ParseIntBasedZipCodeValidator"; //  exports the 'ParseIntBasedZipCodeValidator' class