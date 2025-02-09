#!/bin/bash
python -m venv venv
source venv/bin/activate  # On Windows, use: .\venv\Scripts\activate
pip install -r requirements.txt
pip install -r requirements-dev.txt
cp .env.example .env
