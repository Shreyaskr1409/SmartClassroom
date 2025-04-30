from flask import Flask, request, jsonify
from PIL import Image
import io
import numpy as np
import face_recognition
import pickle
import os

app = Flask(__name__)

# Load encodings and names
MODEL_PATH = "./mlapi/model/face_recognition_model.pkl"
with open(MODEL_PATH, "rb") as f:
    data = pickle.load(f)

known_encodings = data["encodings"]
known_names = data["names"]

@app.route('/predict', methods=['POST'])
def predict():
    if 'image' not in request.files:
        return jsonify({"error": "No image part in the request"}), 400

    file = request.files['image']
    if file.filename == '':
        return jsonify({"error": "No selected image"}), 400

    try:
        image = face_recognition.load_image_file(file)
        face_encodings = face_recognition.face_encodings(image)

        if len(face_encodings) == 0:
            return jsonify({"error": "No face detected."}), 200

        results = []

        for i, encoding in enumerate(face_encodings):
            distances = face_recognition.face_distance(known_encodings, encoding)
            best_match_index = np.argmin(distances)
            best_distance = distances[best_match_index]

            if best_distance < 0.5:  # Adjust threshold as needed
                name = known_names[best_match_index]
                results.append({
                    "face": i + 1,
                    "recognized": True,
                    "name": name,
                    "distance": float(best_distance)
                })
            else:
                results.append({
                    "face": i + 1,
                    "recognized": False,
                    "distance": float(best_distance)
                })

        return jsonify({"results": results})

    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    app.run(port=3300, debug=True)
