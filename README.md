# Morphologically Biased Byte-Pair Encoding

`mbpe-dyn` is a research-focused implementation of byte-pair encoding[^1].
The training results are compatible with the corresponding implementation
in [huggingface/tokenizers](https://github.com/huggingface/tokenizers).

`mbpe-dyn` extends the byte-pair encoding training algorithm as follows. Subword segmentations, which are a direct
result of the trained merge rules, can be aligned to a provided gold segmentation. More specifically, the likelihood of
merge rules that don't interfere with the targeted segmentation is increased. This increase can be tuned via a
hyperparameter.

`mbpe-dyn` was initially created to gain some insights into how well byte-pair encoding approximates morphological
boundaries. Literature suggests that byte-pair encoding produces subword boundaries that align poorly with
linguistically meaningful reference segmentations[^2].

[^1]: [Neural Machine Translation of Rare Words with Subword Units](https://doi.org/10.48550/arXiv.1508.07909)

[^2]: [Byte Pair Encoding is Suboptimal for Language Model Pretraining](https://arxiv.org/abs/2004.03720)

## Related Work

* [MorphPiece : A Linguistic Tokenizer for Large Language Models](https://arxiv.org/abs/2307.07262)
* [MorphBPE: A Morpho-Aware Tokenizer Bridging Linguistic Complexity for Efficient LLM Training Across Morphologies](https://arxiv.org/abs/2502.00894)
* [BPE-knockout: Pruning Pre-existing BPE Tokenisers with Backwards-compatible Morphological Semi-supervision](https://aclanthology.org/2024.naacl-long.324/)

## Limitations

When employing a tokenizer with a close-to-one fertility for large language model training, the intermediate subword
segmentations during tokenization practically do not matter, since they are not relayed to the language model in any
way. Therefore, we suspect our extension to be more useful in settings with tokenization fertility above one.

## Segmenters

| Segmenter   | Description                                                                                                                                                                                                                                                                                                                               |
|-------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `static`    | The `static` segmenter is intended to be used with a morphological lexicon. Lexicon data can be loaded from `.tsv` files that match the format of the [Morpho Challenge](http://morpho.aalto.fi/events/morphochallenge/)[^3] datasets.                                                                                                    |
| `morfessor` | The `morfessor` segmenter allows segmentation via a [Morfessor](https://github.com/aalto-speech/morfessor) baseline model. Trained baseline models need to be converted to a binary format using our Protobuf definition. See the [morfessor](https://github.com/jonasknobloch/mbpe-dyn/tree/main/morfessor)[^4][^5] package for details. |

[^3]: [Morpho Challenge 2005-2010: Evaluations and Results](https://aclanthology.org/W10-2211/)

[^4]: [Unsupervised Discovery of Morphemes](https://doi.org/10.48550/arXiv.cs/0205057)

[^5]: [Morfessor 2.0: Python Implementation and Extensions for Morfessor Baseline](https://urn.fi/URN:ISBN:978-952-60-5501-5)

## Evaluators

| Evaluator                         | Description                                                                                                                                                                                                                                                                                                             |
|-----------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Fertility**                     | Measures the number of tokens per tokenized word. A fertility of one is ideal, meaning each input word maps to exactly one token. Byte-pair encoding optimizes tokenization fertility by prioritizing frequent subwords.                                                                                                |
| **Boundary Precision and Recall** | Evaluates word segmentations by matching them against a gold standard, calculating precision and recall. Based on the [Morpho Challenge](http://morpho.aalto.fi/events/morphochallenge/) evaluation scripts. A partial port is available in the [bpr](https://github.com/jonasknobloch/mbpe-dyn/tree/main/bpr) package. |
| **Merge Layer**                   | Instead of only evaluating the final segmentation, this metric examines intermediate segmentations after each merge. It records the number of previous merges at which a morphological boundary was crossed.                                                                                                            |

## Evaluation Results

![plot](assets/plot.svg)

| #          | Vocabulary | Boundary Precision Recall | Merge Layer | Fertility | Reference Overlap |
|------------|------------|---------------------------|-------------|-----------|-------------------|
| 00-en-m000 | 65536      | 0.65, 0.37, 0.47          | 0.78        | 1.07      | 1.00, 1.00        |
| 01-en-m010 | 65536      | 0.66, 0.38, 0.48          | 0.79        | 1.07      | 0.97, 0.92        |
| 02-en-m020 | 65536      | 0.65, 0.39, 0.49          | 0.79        | 1.07      | 0.94, 0.88        |
| 03-en-m030 | 65536      | 0.65, 0.40, 0.49          | 0.80        | 1.07      | 0.92, 0.84        |
| 04-en-m040 | 65536      | 0.65, 0.40, 0.50          | 0.80        | 1.07      | 0.90, 0.81        |
| 05-en-m050 | 65536      | 0.65, 0.41, 0.50          | 0.80        | 1.07      | 0.88, 0.78        |
| 06-en-m060 | 65536      | 0.64, 0.42, 0.51          | 0.81        | 1.07      | 0.85, 0.75        |
| 07-en-m070 | 65536      | 0.64, 0.43, 0.51          | 0.81        | 1.08      | 0.82, 0.72        |
| 08-en-m080 | 65536      | 0.63, 0.43, 0.51          | 0.81        | 1.08      | 0.78, 0.68        |
| 09-en-m090 | 65536      | 0.63, 0.44, 0.52          | 0.82        | 1.09      | 0.74, 0.64        |
| 10-en-m100 | 65536      | 0.62, 0.47, 0.53          | 0.82        | 1.17      | 0.64, 0.56        |

| #          | Vocabulary | Boundary Precision Recall | Merge Layer | Fertility | Reference Overlap |
|------------|------------|---------------------------|-------------|-----------|-------------------|
| 00-en-m000 | 32768      | 0.54, 0.44, 0.49          | 0.80        | 1.12      | 1.00, 1.00        |
| 01-en-m010 | 32768      | 0.55, 0.45, 0.49          | 0.81        | 1.12      | 0.96, 0.91        |
| 02-en-m020 | 32768      | 0.55, 0.45, 0.50          | 0.81        | 1.12      | 0.94, 0.87        |
| 03-en-m030 | 32768      | 0.56, 0.46, 0.50          | 0.82        | 1.12      | 0.92, 0.84        |
| 04-en-m040 | 32768      | 0.56, 0.46, 0.51          | 0.82        | 1.12      | 0.89, 0.80        |
| 05-en-m050 | 32768      | 0.56, 0.47, 0.51          | 0.82        | 1.12      | 0.87, 0.77        |
| 06-en-m060 | 32768      | 0.56, 0.47, 0.51          | 0.82        | 1.13      | 0.84, 0.74        |
| 07-en-m070 | 32768      | 0.57, 0.48, 0.52          | 0.83        | 1.13      | 0.81, 0.71        |
| 08-en-m080 | 32768      | 0.58, 0.48, 0.52          | 0.83        | 1.14      | 0.78, 0.68        |
| 09-en-m090 | 32768      | 0.57, 0.48, 0.52          | 0.83        | 1.15      | 0.74, 0.64        |
| 10-en-m100 | 32768      | 0.58, 0.49, 0.53          | 0.83        | 1.22      | 0.67, 0.58        |

| #          | Vocabulary | Boundary Precision Recall | Merge Layer | Fertility | Reference Overlap |
|------------|------------|---------------------------|-------------|-----------|-------------------|
| 00-en-m000 | 16384      | 0.48, 0.48, 0.48          | 0.82        | 1.20      | 1.00, 1.00        |
| 01-en-m010 | 16384      | 0.49, 0.49, 0.49          | 0.83        | 1.20      | 0.96, 0.92        |
| 02-en-m020 | 16384      | 0.49, 0.49, 0.49          | 0.83        | 1.20      | 0.93, 0.87        |
| 03-en-m030 | 16384      | 0.49, 0.50, 0.49          | 0.83        | 1.20      | 0.91, 0.84        |
| 04-en-m040 | 16384      | 0.50, 0.50, 0.50          | 0.83        | 1.20      | 0.88, 0.80        |
| 05-en-m050 | 16384      | 0.50, 0.51, 0.50          | 0.83        | 1.21      | 0.86, 0.77        |
| 06-en-m060 | 16384      | 0.51, 0.51, 0.51          | 0.84        | 1.21      | 0.84, 0.75        |
| 07-en-m070 | 16384      | 0.52, 0.52, 0.52          | 0.84        | 1.22      | 0.81, 0.73        |
| 08-en-m080 | 16384      | 0.52, 0.53, 0.52          | 0.84        | 1.22      | 0.78, 0.69        |
| 09-en-m090 | 16384      | 0.52, 0.53, 0.53          | 0.84        | 1.24      | 0.76, 0.67        |
| 10-en-m100 | 16384      | 0.53, 0.53, 0.53          | 0.84        | 1.28      | 0.71, 0.62        |