download-torch:
	wget http://download.pytorch.org/whl/cpu/torch-2.9.0%2Bcpu-cp312-cp312-manylinux_2_28_x86_64.whl#sha256=28f6eb31b08180a5c5e98d5bc14eef6909c9f5a1dbff9632c3e02a8773449349

download-rubert:
	mkdir -p rubert
	cd rubert && wget "https://huggingface.co/DeepPavlov/rubert-base-cased/resolve/main/config.json?download=true" -O config.json
	cd rubert && wget "https://huggingface.co/DeepPavlov/rubert-base-cased/resolve/main/flax_model.msgpack?download=true" -O flax_model.msgpack
	cd rubert && wget "https://huggingface.co/DeepPavlov/rubert-base-cased/resolve/main/pytorch_model.bin?download=true" -O pytorch_model.bin
	cd rubert && wget "https://huggingface.co/DeepPavlov/rubert-base-cased/resolve/main/special_tokens_map.json?download=true" -O special_tokens_map.json
	cd rubert && wget "https://huggingface.co/DeepPavlov/rubert-base-cased/resolve/main/tokenizer_config.json?download=true" -O tokenizer_config.json
	cd rubert && wget "https://huggingface.co/DeepPavlov/rubert-base-cased/resolve/main/vocab.txt?download=true" -O vocab.txt

download: download-torch download-rubert
