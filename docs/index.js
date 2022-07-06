import { renderDebugInfo } from "./debugger/index.js";

window.onload = async () => {
  document.addEventListener(
    "touchmove",
    function(event) {
      if (event.scale !== 1) {
        event.preventDefault();
      }
    },
    false
  );

  const go = new Go();
  const res = await fetch("./main.wasm");
  const bytes = await res.arrayBuffer();
  const { instance } = await WebAssembly.instantiate(bytes, go.importObject);
  go.run(instance);

  const canvas = document.querySelector(".game");
  const input = document.querySelector("#file_upload");

  // GPU

  const init = async buf => {
    const ctx = canvas.getContext("2d");
    const image = ctx.createImageData(160, 144);

    if (!buf) {
      const rom = await fetch("./tobu.gb");
      buf = await rom.arrayBuffer();
    }
    let gb = new GB(new Uint8Array(buf));

    document.querySelector(".led").style.background = "red";

    const frame = () => {
      if (!gb) return;
      gb.next(image.data);
      ctx.putImageData(image, 0, 0);
      renderDebugInfo(gb);
      window.requestAnimationFrame(frame);
    };
    frame();

    const onKeydown = e => {
      switch (e.key) {
        case "z":
          return gb.keyDown(0x01);
        case "x":
          return gb.keyDown(0x02);
        case "Backspace":
          return gb.keyDown(0x04);
        case "Enter":
          return gb.keyDown(0x08);
        case "ArrowLeft":
          return gb.keyDown(0x20);
        case "ArrowUp":
          return gb.keyDown(0x40);
        case "ArrowRight":
          return gb.keyDown(0x10);
        case "ArrowDown":
          return gb.keyDown(0x80);
      }
    };

    const onKeyup = e => {
      switch (e.key) {
        case "z":
          return gb.keyUp(0x01);
        case "x":
          return gb.keyUp(0x02);
        case "Backspace":
          return gb.keyUp(0x04);
        case "Enter":
          return gb.keyUp(0x08);
        case "ArrowLeft":
          return gb.keyUp(0x20);
        case "ArrowUp":
          return gb.keyUp(0x40);
        case "ArrowRight":
          return gb.keyUp(0x10);
        case "ArrowDown":
          return gb.keyUp(0x80);
      }
    };

    const removeHandler = classname => {
      const el = document.querySelector(`.${classname}`);
      const elClone = el.cloneNode(true);
      el.parentNode.replaceChild(elClone, el);
    };

    const cleanup = () => {
      input.removeEventListener("change", onFileChange);
      window.removeEventListener("keydown", onKeydown);
      window.removeEventListener("keyup", onKeyup);
      removeHandler("buttonA");
      removeHandler("buttonB");