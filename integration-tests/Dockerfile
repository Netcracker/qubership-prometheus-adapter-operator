FROM python:3.9-slim

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

WORKDIR /app

COPY python-scripts/check-metric-apiservice.py .
CMD ["python", "check-metric-apiservice.py"]