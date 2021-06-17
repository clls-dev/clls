import 'regenerator-runtime/runtime'
import "./wasm_exec"
import binaryDataUrl from 'data-url:./main.wasm';

function mapTypeToColor(t, m) {
    if (t == 12 && m != 0) { // builtin func
        return "rgb(255, 255, 145)"
    }
    return ["red", "green", "blue", "purple", "pink", "grey", "yellow", "rgb(144, 238, 144)", "rgb(185, 185, 255)", "grey", "purple", "pink", "rgb(255, 255, 65)", "yellow", "red", "rgb(200, 0, 200)", "blue", "rgb(206, 208, 210)", "rgb(201, 139, 50)", "white", "yellow", "red", "green", "blue", "purple", "pink", "grey", "yellow"][t]
}

function color(toks, input) {
    const inserts = []

    const defColor = "rgb(125, 125, 255)"
    const constColor = "rgb(185, 185, 255)"
    const keywordColor = "rgb(200, 0, 200)"
    const funcColor = "rgb(255, 255, 65)"
    const builtinColor = "rgb(255, 255, 145)"
    const argColor = "rgb(144, 238, 144)"
    const templateColor = "rgb(255, 165, 0)"
    const commentColor = "rgb(206, 208, 210)"
    const literalColor = "rgb(201, 139, 50)"

    for (let i = 0; i * 5 < toks.length; i++) {
        const [deltaLine, deltaStartChar, len, tokenType, tokenModifiers] = toks.slice(i * 5, (i + 1) * 5)
        const color = mapTypeToColor(tokenType, tokenModifiers)
        inserts.push({ color, deltaLine, deltaStartChar, len })
    }

    const lines = input.split("\n")
    let final = ""
    let lineIndex = 0
    let lineRemains = lines[lineIndex]
    prevLen = 0
    for (const ins of inserts) {
        if (ins.deltaLine > 0) {
            final += lineRemains + "\n"
            for (let i = 1; i < ins.deltaLine; i++) {
                final += lines[lineIndex + i] + "\n"
            }
            lineIndex += ins.deltaLine
            lineRemains = lines[lineIndex]
            prevLen = 0
        }
        final += lineRemains.substr(0, ins.deltaStartChar - prevLen)
        lineRemains = lineRemains.substr(ins.deltaStartChar - prevLen)
        final += '<span style="color: ' + ins.color + '">' +
            lineRemains.substr(0, ins.len) +
            "</span>"
        lineRemains = lineRemains.substr(ins.len)
        prevLen = ins.len
    }
    final += lineRemains + "\n"
    lineIndex++
    for (; lineIndex < lines.length; lineIndex++) {
        final += lines[lineIndex]
        lineIndex++
        if (lineIndex < lines.length) {
            final += "\n"
        }
    }
    return final
}

(async function amain() {
    const go = new Go();
    const buffer = await (await fetch(binaryDataUrl)).arrayBuffer()
    const result = await WebAssembly.instantiate(buffer, go.importObject);
    go.run(result.instance)

    let running = false
    let rerun = true

    async function update(input) {
        if (running) {
            rerun = true
            return
        }
        running = true

        const toks = window.cllsSemanticTokens(input)
        if (toks) {
            if (Array.isArray(toks)) {
                document.getElementById("main").innerHTML = color(toks, input)
            } else if (typeof out == "string") {
                document.getElementById("main").innerHTML = toks
            } else {
                console.log(toks)
            }
        }
        running = false
        if (rerun) {
            rerun = false
            const ine = document.getElementById("intext");
            update(ine.value)
        }
    }



    const examples = window.cllsGetExamples()
    document.getElementById("examples").innerHTML = '<select id="examples-select">' + Object.entries(examples).sort(([a], [b]) => a > b).reduce((result, [exampleName, exampleCode]) => {
        return result + '"<option value="' + exampleName + '">' + exampleName + '</option>'
    }, "") + "</select>"

    document.getElementById("examples-select").value = "p2_delegated_puzzle_or_hidden_puzzle.clvm"
    document.getElementById("intext").value = examples["p2_delegated_puzzle_or_hidden_puzzle.clvm"]
    update(examples["p2_delegated_puzzle_or_hidden_puzzle.clvm"])

    document.getElementById("examples-select").addEventListener("change", (e) => {
        document.getElementById("intext").value = examples[e.target.value]
        update(examples[e.target.value])
    })
    document.getElementById("intext").addEventListener("input", (e) => {
        update(e.target.value)
    })
})()