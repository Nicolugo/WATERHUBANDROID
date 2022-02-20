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
const gpuOBP1 =