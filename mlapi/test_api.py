import unittest
import os
import requests
import time
import json
import subprocess
from PIL import Image

class MLAPITest(unittest.TestCase):
    API_URL = "http://localhost:3300/predict"
    SERVER_PROCESS = None

    @classmethod
    def setUpClass(cls):
        # Start the Flask server for testing
        print("Starting Flask server for testing...")
        cls.SERVER_PROCESS = subprocess.Popen(
            ["python3", "app.py"], 
            cwd="./mlapi",
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE
        )
        # Give the server time to start
        time.sleep(3)
        print("Flask server started")

    @classmethod
    def tearDownClass(cls):
        # Shut down the server after all tests
        if cls.SERVER_PROCESS:
            cls.SERVER_PROCESS.terminate()
            cls.SERVER_PROCESS.wait()
            print("Flask server stopped")

    def test_server_running(self):
        """Test if the API server is running"""
        try:
            response = requests.get("http://localhost:3300/")
            self.assertEqual(response.status_code, 404)  # Flask returns 404 for undefined routes
            print("API server is running")
        except requests.exceptions.ConnectionError:
            self.fail("API server is not running")

    def test_predict_endpoint_no_image(self):
        """Test the predict endpoint with no image"""
        response = requests.post(self.API_URL)
        self.assertEqual(response.status_code, 400)
        self.assertIn("error", response.json())
        print("Predict endpoint correctly handles missing image")

    def test_predict_endpoint_with_image(self):
        """Test the predict endpoint with a test image"""
        # Create a simple test image
        test_image_path = "/tmp/test_face.jpg"
        
        try:
            # Check if the test image exists, otherwise skip
            if not os.path.exists(test_image_path):
                print(f"Creating test image at {test_image_path}")
                img = Image.new('RGB', (640, 480), color=(255, 255, 255))
                img.save(test_image_path)
            
            # Send the image to the API
            with open(test_image_path, 'rb') as img_file:
                files = {'image': img_file}
                response = requests.post(self.API_URL, files=files)
            
            # Check the response
            self.assertEqual(response.status_code, 200)
            
            # The response might indicate no face was found, which is expected with a blank test image
            # But the format should be correct
            response_data = response.json()
            if "error" in response_data and response_data["error"] == "No face detected.":
                print("API correctly reported no face in test image")
            elif "results" in response_data:
                print(f"API processed image and returned results: {response_data}")
            else:
                self.fail(f"Unexpected API response: {response_data}")
                
        finally:
            # Clean up the test image
            if os.path.exists(test_image_path):
                os.remove(test_image_path)
                print(f"Removed test image {test_image_path}")

if __name__ == "__main__":
    unittest.main()
