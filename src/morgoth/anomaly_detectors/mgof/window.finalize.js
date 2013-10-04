function (key, value) {
    var len = value.count * 10 + value.P.length;
    for (var i = 0; i < value.P.length; i++) {
        value.P[i] /= len;
    }
    delete value['discrete'];
    return value;
}
