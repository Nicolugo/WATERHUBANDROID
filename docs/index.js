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

    if (!buf