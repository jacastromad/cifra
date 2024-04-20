#!/bin/bash

# Integration test - cifra
#
# Build cifra and test it against randomly generated files

fname="itest.dat"

text_file() {
    cat /dev/urandom | tr -dc '[:print:]\n' | head -n 1000 > "$fname"
}

binary_file() {
    dd if=/dev/urandom of="$fname" bs=10000 count=1 status=none
}

create_files=(text_file binary_file)


#########
# Tests #
#########

test1() {
    echo "1234" | ./cifra "$fname" > /dev/null && \
    echo "1234" | ./cifra -dec "${fname}.cif" > /dev/null && \
    cmp -s "${fname}.cif.dec" "${fname}"
}

test2() {
    echo "1234" | ./cifra -cfb -o "${fname}.enc" "$fname" > /dev/null && \
    echo "1234" | ./cifra -dec -cfb -o "${fname}.dec" "${fname}.enc" > /dev/null && \
    cmp -s "${fname}.dec" "${fname}"
}

test3() {
    echo "1234" | ./cifra -b64 -o "itest.b64.enc" "$fname" > /dev/null && \
    echo "1234" | ./cifra -b64 -dec -o "itest.b64.dec" "itest.b64.enc" > /dev/null && \
    cmp -s "itest.b64.dec" "${fname}"
}

tests=(test1 test2 test3)

# Build cifra for integration test
go build -tags=itest -o cifra
if [ $? -ne 0 ]; then
    echo "Failed to build cifra for integration test"
    exit 1
fi

#############
# Run tests #
#############

for create_file in "${create_files[@]}"; do
    echo "Creating ${create_file}..."
    $create_file || {
        echo "Error creating ${create_file}"
        exit 1
    }
    for tst in "${tests[@]}"; do
        $tst && echo "${tst} ok" || {
            echo "Test ${tst} failed for ${create_file}"
            rm itest.* cifra
            exit 1
        }
    done
    # clean up
    rm itest.*
done

# Build cifra
go build -o cifra
if [ $? -ne 0 ]; then
    echo "Failed to build cifra"
    exit 1
fi

echo "All tests passed successfully."

