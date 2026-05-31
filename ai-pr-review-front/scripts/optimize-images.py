"""将 png/ 原图压缩为 WebP，供前端静态资源使用。"""
from pathlib import Path

from PIL import Image

ROOT = Path(__file__).resolve().parent.parent
SRC_DIR = ROOT / "png"
DST_DIR = ROOT / "public" / "images" / "decor"

# (源文件, 输出名, 最大宽度, 质量)
IMAGES = [
    ("144194348_p0.jpg", "washitsu.webp", 1600, 72),
    ("144167969_p0.jpg", "warrior.webp", 900, 75),
    ("144183490_p0.png", "pool.webp", 800, 80),
    ("144158180_p0.jpg", "campfire.webp", 1000, 75),
    ("144153210_p0.jpg", "catgirl.webp", 700, 80),
    ("2.jpg", "history-bg.webp", 1600, 78),
]


def optimize(src_name: str, dst_name: str, max_width: int, quality: int) -> None:
    src = SRC_DIR / src_name
    dst = DST_DIR / dst_name
    if not src.exists():
        print(f"skip missing: {src}")
        return

    img = Image.open(src).convert("RGB")
    w, h = img.size
    if w > max_width:
        ratio = max_width / w
        img = img.resize((max_width, int(h * ratio)), Image.Resampling.LANCZOS)

    DST_DIR.mkdir(parents=True, exist_ok=True)
    img.save(dst, "WEBP", quality=quality, method=6)
    print(f"{src_name} -> {dst_name} ({dst.stat().st_size // 1024} KB)")


if __name__ == "__main__":
    for item in IMAGES:
        optimize(*item)
