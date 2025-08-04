import pyttsx3

def text_to_speech_converter():
    """
    Converts text input from the user to speech using a platform-specific TTS engine.
    """
    engine = pyttsx3.init()

    # Optional: Configure voice properties
    # voices = engine.getProperty('voices')
    # for voice in voices:
    #     print(f"Voice ID: {voice.id}, Name: {voice.name}, Gender: {voice.gender}, Age: {voice.age}")
    # engine.setProperty('voice', voices[0].id) # Set to a specific voice, e.g., the first available

    rate = engine.getProperty('rate')
    engine.setProperty('rate', 150) # Speed of speech (words per minute)

    volume = engine.getProperty('volume')
    engine.setProperty('volume', 1.0) # Volume (0.0 to 1.0)

    print("Text-to-Speech Converter (Type 'exit' to quit)")
    while True:
        text_input = input("Enter text to convert to speech: ")
        if text_input.lower() == 'exit':
            break
        
        engine.say(text_input)
        engine.runAndWait()

    engine.stop()

if __name__ == "__main__":
    text_to_speech_converter()

# Additional implementation at 2025-08-04 07:15:00
import pyttsx3
import os

class TTSConverter:
    def __init__(self):
        self.engine = pyttsx3.init()
        self._set_default_properties()

    def _set_default_properties(self):
        self.engine.setProperty('rate', 150)
        self.engine.setProperty('volume', 0.9)

    def speak(self, text):
        self.engine.say(text)
        self.engine.runAndWait()

    def save_to_file(self, text, filename="output.mp3"):
        try:
            self.engine.save_to_file(text, filename)
            self.engine.runAndWait()
            print(f"Text saved to {filename}")
        except Exception as e:
            print(f"Error saving to file: {e}")

    def set_rate(self, rate):
        if isinstance(rate, int) and rate > 0:
            self.engine.setProperty('rate', rate)
        else:
            print("Invalid rate. Please provide a positive integer.")

    def set_volume(self, volume):
        if isinstance(volume, (int, float)) and 0.0 <= volume <= 1.0:
            self.engine.setProperty('volume', volume)
        else:
            print("Invalid volume. Please provide a value between 0.0 and 1.0.")

    def list_voices(self):
        voices = self.engine.getProperty('voices')
        print("Available Voices:")
        for i, voice in enumerate(voices):
            print(f"{i+1}. ID: {voice.id}, Name: {voice.name}, Lang: {voice.languages[0] if voice.languages else 'N/A'}, Gender: {voice.gender if voice.gender else 'N/A'}")
        return voices

    def set_voice(self, voice_id):
        voices = self.engine.getProperty('voices')
        found = False
        for voice in voices:
            if voice.id == voice_id:
                self.engine.setProperty('voice', voice.id)
                found = True
                break
        if not found:
            print(f"Voice with ID '{voice_id}' not found. Using default voice.")

    def read_text_from_file(self, filepath):
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                text = f.read()
            return text
        except FileNotFoundError:
            print(f"Error: File not found at {filepath}")
            return None
        except Exception as e:
            print(f"Error reading file: {e}")
            return None

    def stop(self):
        self.engine.stop()

if __name__ == "__main__":
    tts = TTSConverter()

    tts.speak("Hello, this is a basic text-to-speech example.")

    tts.set_rate(180)
    tts.set_volume(0.7)
    tts.speak("I am now speaking a bit faster and quieter.")

    available_voices = tts.list_voices()
    if available_voices:
        selected_voice_id = None
        for voice in available_voices:
            if "female" in voice.gender.lower():
                selected_voice_id = voice.id
                break
        if not selected_voice_id:
            selected_voice_id = available_voices[0].id

        tts.set_voice(selected_voice_id)
        tts.speak("Now I am speaking with a different voice.")
    else:
        tts.speak("No alternative voices found.")

    tts.save_to_file("This text will be saved to an MP3 file.", "saved_audio.mp3")

    dummy_file_path = "sample_text.txt"
    with open(dummy_file_path, "w", encoding="utf-8") as f:
        f.write("This is some text read from a file. It demonstrates reading capabilities.")

    file_content = tts.read_text_from_file(dummy_file_path)
    if file_content:
        tts.speak(file_content)
        tts.save_to_file(file_content, "file_audio.mp3")

    if os.path.exists(dummy_file_path):
        os.remove(dummy_file_path)

    tts.stop()