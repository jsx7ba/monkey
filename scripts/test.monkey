let makeHash = fn(a, b, c) {
    let h = {1: a, 2: b, 3: c};
    return h
};

let a = "1";
let b = "2";
let c = "3";

let h = makeHash(a, b, c);
puts(h)
