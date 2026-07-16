from faster_whisper import WhisperModel
from pathlib import Path

MODEL_SIZE = "small"

model = WhisperModel(
    MODEL_SIZE,
    device="cpu",
    compute_type="int8"
)

def transcribe(audio_path: str):
    segments, info = model.transcribe(
        audio_path,
        beam_size=1,
        vad_filter=True
    )

    print(f"Language: {info.language}")
    print(f"Confidence: {info.language_probability:.2f}")

    transcript = []

    for segment in segments:
        line = (
            f"[{segment.start:.2f}s -> "
            f"{segment.end:.2f}s] "
            f"{segment.text}"
        )

        print(line)
        transcript.append(line)

    return "\n".join(transcript)


if __name__ == "__main__":
    audio_file = "download/audio.mp3"

    text = transcribe(audio_file)

    Path("output/transcript.txt").write_text(
        text,
        encoding="utf-8"
    )

    print("\nSaved to transcript.txt")
