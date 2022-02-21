const getPalette = c => {
  switch (c) {
    case 0:
      return [175, 197, 160, 255];
    case 1:
      return [93, 147, 66, 255];
    case 2:
      return [22, 63, 48, 255];
    case 3:
      return [0, 40, 0, 255];
  }
};

const gpuLCDC = document.querySelector(".gpu-lcdc");
const gpuSTAT = document.querySelector(".gpu-stat");
const gpuSCROLLY = document.querySelector(".gpu-scrolly");
const gpuSCROLLX = document.querySelector(".gpu-scrollx");
const gpuLY = document.querySelector(".gpu-ly");
const gpuLYC = document.querySelector(".gpu-lyc");
const gpuDMA = document.querySelector(".gpu-dma");
const gpuBGP = document.querySelector(".gpu-bgp");
const gpuOBP0 = document.querySelector(".gpu-obp0");
const gpuOBP1 = document.querySelector(".gpu-obp1");
const gpuWY = document.querySelector(".gpu-wy");
const gpuWX = document.querySelector(".gpu-wx");

const tileData0 = document.querySelector(".tiledata0");

const renderTileMap = (ctx, vram, tiles, offsetAddr, tileData0Selected) => {
  const tileMap = [];
  for (let n = 0; n < 1024; n++) {
    const tileId = vram[offsetAddr + n];
    let index = tileId;
    if (tileData0Selected) {
      index =
        (tileId & 0x80 ? new Int8Array([tileId])[0] : tileId & 0x7f) + 256;
    }
    const sprite = tiles[index];
    for (let i = 0; i < 8; i++) {
      for (let j = 0; j < 8; j++) {
        const c = getPalette(sprite[i][j]);
        const x = j + (n % 32) * 8;
        const y = i + ~~(n / 32) * 8;
        tileMap[(y * 256 + x) * 4] = c[0];
        tileMap[(y * 256 + x) * 4 + 1] = c[1];
        tileMap[(y * 256 + x) * 4 + 2] = c[2];
        tileMap[(y * 256 + x) * 4 + 3] = 255;
      }
    }
  }
  const image = ctx.createImageData(256, 256);
  image.data.set(tileMap);
  ctx.putImageData(image, 0, 0);
};

const renderWindow = (ctx, vram, tiles, offsetAddr, wx, wy) => {
  const windowMap = [];
  for (let n = 0; n < 640; n++) {
    const tileId = vram[offsetAddr + n];
    let index = tileId;
    index = (tileId & 0x80 ? new Int8Array([tileId])[0] : tileId & 0x7f) + 256;
    const sprite = tiles[index];
    