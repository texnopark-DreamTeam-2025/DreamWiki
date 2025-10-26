FROM python:3.12

WORKDIR /app

RUN apt-get update

RUN apt-get install -y \
    gcc wget \
    && rm -rf /var/lib/apt/lists/*

COPY torch-2.9.0+cpu-cp312-cp312-manylinux_2_28_x86_64.whl .
RUN pip3 install torch-2.9.0+cpu-cp312-cp312-manylinux_2_28_x86_64.whl
RUN rm torch-2.9.0+cpu-cp312-cp312-manylinux_2_28_x86_64.whl

COPY requirements.txt .
RUN pip3 install --no-cache-dir -r requirements.txt

COPY . .

EXPOSE 8000

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
