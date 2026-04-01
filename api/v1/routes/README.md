"""
api-service
-----------

A RESTful API service using Flask.

Usage
-----

To run the service, use the following command:

    python app.py

This will start the service and make it available on ``http://localhost:5000``.
"""

import os

from flask import Flask
from flask_sqlalchemy import SQLAlchemy

# Create the Flask application
app = Flask(__name__)

# Configure the database connection
app.config['SQLALCHEMY_DATABASE_URI'] = os.environ.get(
    'DATABASE_URL',
    'sqlite:///api-service.db'
)
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False

# Create the database instance
db = SQLAlchemy(app)

# Import the routes
from . import routes