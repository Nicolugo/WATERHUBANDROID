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
          return gb.keyDown