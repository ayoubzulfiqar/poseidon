import pyttsx3

def text_to_speech_converter():
    """
    Converts text input from the user to speech using the default
    platform-specific text-to-speech engine.
    """
    try:
        engine = pyttsx3.init()
    except Exception as e:
        print(f"Error initializing text-to-speech engine: {e}")
        print("Please ensure you have a compatible TTS engine installed on your system.")
        print("For Windows: SAPI5 is usually built-in.")
        print("For macOS: NSSpeechSynthesizer is usually built-in.")
        print("For Linux: You might need to install 'espeak' or 'festival'.")
        print("  e.g., sudo apt-get install espeak")
        return

    # Optional: Configure properties (voice, rate, volume)
    # You can uncomment and modify these lines to customize the speech
    # voices = engine.getProperty('voices')
    # if voices:
    #     # Try to set a female voice if available, otherwise use default
    #     found_female = False
    #     for voice in voices:
    #         # Check for common indicators of female voices
    #         if 'female' in voice.name.lower() or 'zira' in voice.name.lower() or 'helen' in voice.name.lower():
    #             engine.setProperty('voice', voice.id)
    #             found_female = True
    #             break
    #     if not found_female and voices:
    #         engine.setProperty('voice', voices[0].id) # Default to first available voice
    # engine.setProperty('rate', 175)    # Speed of speech (words per minute)
    # engine.setProperty('volume', 1.0)  # Volume (0.0 to 1.0)

    print("Text-to-Speech Converter")
    print("Enter the text you want to convert to speech.")
    print("Type 'exit' or 'quit' to stop the program.")

    while True:
        try:
            text_input = input("Your text: ")
            if text_input.lower() in ['exit', 'quit']:
                break
            if text_input.strip():
                engine.say(text_input)
                engine.runAndWait()
            else:
                print("No text entered. Please type something or 'exit'.")
        except Exception as e:
            print(f"An error occurred during speech generation: {e}")
            break # Exit loop on critical error

    engine.stop()
    print("Text-to-speech converter stopped.")

if __name__ == "__main__":
    text_to_speech_converter()

# Additional implementation at 2025-08-04 07:45:08
import pyttsx3
import os

class TextToSpeechConverter:
    def __init__(self, voice_id=None, rate=150, volume=1.0):
        self.engine = pyttsx3.init()
        self.set_properties(rate, volume)
        
        # Attempt to set a specific voice if provided, otherwise try to pick a default
        if voice_id:
            self.set_voice(voice_id)
        else:
            self._set_default_voice()

    def _set_default_voice(self):
        voices = self.engine.getProperty('voices')
        if voices:
            # Prioritize an English female voice if available, otherwise any English, then first available
            selected_voice_id = None
            for voice in voices:
                # Check for common indicators of English and female voices
                if ("en" in voice.languages or "english" in voice.name.lower()) and \
                   ("female" in voice.name.lower() or "f" == getattr(voice, 'gender', '').lower()):
                    selected_voice_id = voice.id
                    break
            
            if not selected_voice_id:
                for voice in voices:
                    if "en" in voice.languages or "english" in voice.name.lower():
                        selected_voice_id = voice.id
                        break
            
            if not selected_voice_id:
                selected_voice_id = voices[0].id # Fallback to the very first voice

            if selected_voice_id:
                self.set_voice(selected_voice_id)
            else:
                print("Warning: No voices found on the system.")
        else:
            print("Warning: No voices available to set.")

    def set_properties(self, rate=None, volume=None):
        if rate is not None:
            self.engine.setProperty('rate', rate)
        if volume is not None:
            self.engine.setProperty('volume', volume)

    def set_voice(self, voice_id):
        voices = self.engine.getProperty('voices')
        found = False
        for voice in voices:
            if voice.id == voice_id:
                self.engine.setProperty('voice', voice_id)
                found = True
                break
        if not found:
            print(f"Warning: Voice ID '{voice_id}' not found. Using current voice.")
        return found

    def list_voices(self):
        voices = self.engine.getProperty('voices')
        print("--- Available Voices ---")
        if not voices:
            print("No voices found on this system.")
            return []
        
        for i, voice in enumerate(voices):
            print(f"  {i+1}. ID: {voice.id}")
            print(f"     Name: {voice.name}")
            print(f"     Languages: {voice.languages}")
            print(f"     Gender: {getattr(voice, 'gender', 'N/A')}") # gender might not be present on all drivers
            print(f"     Age: {getattr(voice, 'age', 'N/A')}") # age might not be present on all drivers
            print("-" * 25)
        print("--- End Voice List ---")
        return voices

    def speak(self, text):
        print(f"Speaking: '{text}'")
        self.engine.say(text)
        self.engine.runAndWait()

    def save_to_file(self, text, filename="output.mp3"):
        print(f"Attempting to save to '{filename}'...")
        try:
            self.engine.save_to_file(text, filename)
            self.engine.runAndWait()
            if os.path.exists(filename) and os.path.getsize(filename) > 0:
                print(f"Text successfully saved to {filename}")
            else:
                print(f"Warning: File '{filename}' might not have been created or is empty. Check system dependencies (e.g., ffmpeg for MP3).")
        except Exception as e:
            print(f"Error saving file: {e}")

    def stop(self):
        self.engine.stop()

    def __del__(self):
        if hasattr(self, 'engine'):
            self.engine.stop()

if __name__ == "__main__":
    converter = TextToSpeechConverter(rate=170, volume=0.9)

    print("\n--- Initializing and Listing Voices ---")
    available_voices = converter.list_voices()

    # Example of selecting a specific voice by ID (uncomment and replace with an actual ID from list_voices)
    # if available_voices:
    #     # Example: Try to pick the first voice listed
    #     first_voice_id = available_voices[0].id
    #     print(f"\nAttempting to set voice to the first available: {first_voice_id}")
    #     converter.set_voice(first_voice_id)

    print("\n--- Performing Speech Operations ---")
    
    text1 = "Hello, this is a versatile text-to-speech converter."
    converter.speak(text1)

    converter.set_properties(rate=120, volume=0.7)
    text2 = "I can speak slower and quieter now, demonstrating property changes."
    converter.speak(text2)

    converter.set_properties(rate=200, volume=1.0)
    text3 = "And I can speak much faster and louder if needed!"
    converter.speak(text3)

    print("\n--- Saving Speech to File ---")
    file_text = "This sentence will be saved as an audio file. You can find it in the current directory."
    output_audio_file = "my_spoken_text.mp3"
    converter.save_to_file(file_text, output_audio_file)

    print("\n--- Text-to-speech operations completed. ---")
    # The converter will automatically stop its engine when the script exits or the object is garbage collected.