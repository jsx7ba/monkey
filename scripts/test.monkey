    let foo = fn(a, al) {
        if (al != 0) {
            return foo(push(a, "foobar"), al - 1);
        }
        return a;
    };

let makeArray = fn(size) {
    let r = [];
    return foo(r, size);
};

let length = 2;
let pi = 3.141592;
let array_1 = [pi, 0xDEADBEEF, "test"];
let xx = array_1[1];
puts(xx);
puts(array_1[0]);

let array = makeArray(10);
puts("first element: " + array[0]);
puts(len(array))

let makeHash = fn(a, b, c) {
    let h = {1: a, 2: b, 3: c};
    return h
};

let a = "1";
let b = "2";
let c = "3";

let h = makeHash(a, b, c);

let foo = ["one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"];

let label = fn(foo) {
    let x = first(foo);
    label(rest(foo));
}

label(foo);

let pi = 3.141592653589793;
let phi = 1.618033988749895;
let weird = pi / phi;

puts(weird);
