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
    let x = first(foo) + "-er"
    label(rest)
}

label(foo)

if cmp(1.2001, 1.201) == -1 {
let foo = cmp("2.2", "two-point-two")
}

