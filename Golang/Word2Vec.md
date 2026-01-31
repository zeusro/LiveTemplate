# Word2Vec 原理简述（Go 示例）

Word2Vec 是 Mikolov 等人提出的词嵌入（Word Embedding）方法，将词语映射到低维稠密向量，使得语义相近的词在向量空间中距离更近。

## 核心思想

- **分布式假设**：出现在相似上下文中的词，语义往往相近。
- **目标**：学习一个映射 \( f: \text{word} \to \mathbb{R}^d \)，使得上下文相似的词向量接近。

## 两种架构

| 架构 | 输入 | 输出 | 适用场景 |
|------|------|------|----------|
| **CBOW** (Continuous Bag of Words) | 上下文词（多个） | 中心词 | 语料较少时 |
| **Skip-gram** | 中心词 | 上下文词 | 语料较少时效果更好，常用 |

- **CBOW**：用上下文词的平均向量预测中心词。
- **Skip-gram**：用中心词预测其上下文词（多分类/二分类）。

下面用 Go 写一个**极简的 Skip-gram 风格**示例，只体现「词→向量→预测上下文」的流程，不涉及真实训练。

---

## 1. 词表与 One-Hot 编码

词先映射为 ID，再转为 one-hot 向量（仅用于理解，实际 Word2Vec 里会用嵌入矩阵直接查表）。

```go
package main

import "fmt"

// 词表：词 -> ID
var vocab = map[string]int{
	"我": 0, "爱": 1, "自然": 2, "语言": 3, "处理": 4,
}

// oneHot 生成维度为 vocabSize 的 one-hot 向量
func oneHot(word string, vocabSize int) []float64 {
	vec := make([]float64, vocabSize)
	if id, ok := vocab[word]; ok {
		vec[id] = 1.0
	}
	return vec
}

func main() {
	for word := range vocab {
		fmt.Printf("%s -> %v\n", word, oneHot(word, len(vocab)))
	}
}
```

---

## 2. 嵌入层：词 ID → 稠密向量

Word2Vec 的核心是**嵌入矩阵** \( W \in \mathbb{R}^{V \times d} \)：\( V \) 为词表大小，\( d \) 为嵌入维度。  
词 \( i \) 的向量即 \( W \) 的第 \( i \) 行（或列，视实现而定）。下面用 Go 表示「查表得到向量」。

```go
package main

import "fmt"

var vocab = map[string]int{"我": 0, "爱": 1, "自然": 2, "语言": 3, "处理": 4}

const (
	vocabSize = 5
	embedDim  = 3
)

// 嵌入矩阵 [vocabSize][embedDim]，实际训练中由梯度下降学习
var embedding = [vocabSize][embedDim]float64{
	{0.1, 0.2, 0.1},   // 我
	{0.3, 0.1, 0.4},   // 爱
	{0.2, 0.5, 0.2},   // 自然
	{0.4, 0.1, 0.3},   // 语言
	{0.1, 0.3, 0.5},   // 处理
}

// getEmbedding 根据词取出其嵌入向量
func getEmbedding(word string) [embedDim]float64 {
	id, ok := vocab[word]
	if !ok {
		return [embedDim]float64{}
	}
	return embedding[id]
}

func main() {
	word := "爱"
	vec := getEmbedding(word)
	fmt.Printf("词 \"%s\" 的嵌入向量: %v\n", word, vec)
}
```

---

## 3. Skip-gram 的「中心词 → 上下文」得分

Skip-gram 用中心词向量 \( v_c \) 与每个候选词（如上下文词）的向量 \( u_w \) 做内积，得到未归一化得分，再通过 softmax 得到概率：

\[
P(w_o \mid w_c) = \frac{\exp(u_{w_o}^\top v_c)}{\sum_k \exp(u_k^\top v_c)}
\]

下面用 Go 实现：**中心词 → 与所有词的得分（logits）→ softmax 成概率**。

```go
package main

import (
	"fmt"
	"math"
)

var vocab = map[string]int{"我": 0, "爱": 1, "自然": 2, "语言": 3, "处理": 4}
var idToWord = []string{"我", "爱", "自然", "语言", "处理"}

const vocabSize, embedDim = 5, 3

// 中心词嵌入 (输入嵌入)
var centerEmbed = [vocabSize][embedDim]float64{
	{0.1, 0.2, 0.1}, {0.3, 0.1, 0.4}, {0.2, 0.5, 0.2},
	{0.4, 0.1, 0.3}, {0.1, 0.3, 0.5},
}

// 上下文词嵌入 (输出嵌入)，实际与 centerEmbed 可共享或分开
var contextEmbed = [vocabSize][embedDim]float64{
	{0.15, 0.25, 0.12}, {0.28, 0.12, 0.38}, {0.22, 0.48, 0.18},
	{0.38, 0.15, 0.28}, {0.12, 0.28, 0.52},
}

// dot 向量内积
func dot(a, b [embedDim]float64) float64 {
	var sum float64
	for i := 0; i < embedDim; i++ {
		sum += a[i] * b[i]
	}
	return sum
}

// softmax 将 logits 转为概率
func softmax(logits []float64) []float64 {
	out := make([]float64, len(logits))
	max := logits[0]
	for _, x := range logits[1:] {
		if x > max {
			max = x
		}
	}
	var sum float64
	for i, x := range logits {
		out[i] = math.Exp(x - max)
		sum += out[i]
	}
	for i := range out {
		out[i] /= sum
	}
	return out
}

func main() {
	centerWord := "爱"
	centerID := vocab[centerWord]
	vc := centerEmbed[centerID]

	// 计算与每个词作为上下文时的得分 (logits)
	logits := make([]float64, vocabSize)
	for i := 0; i < vocabSize; i++ {
		logits[i] = dot(vc, contextEmbed[i])
	}

	probs := softmax(logits)
	fmt.Printf("中心词 \"%s\" 预测各词为上下文的概率:\n", centerWord)
	for i, p := range probs {
		fmt.Printf("  %s: %.4f\n", idToWord[i], p)
	}
}
```

---

## 4. 负采样（Negative Sampling）思想

完整 softmax 在词表很大时计算昂贵。负采样用二分类近似：  
- 正样本：真实上下文词 \((w_c, w_o)\)，标签 1；  
- 负样本：随机采样的非上下文词 \((w_c, w_{\text{neg}})\)，标签 0。  

损失常用二元交叉熵。下面只演示「正样本 + 若干负样本」的构造方式（不写完整训练）。

```go
package main

import (
	"fmt"
	"math"
	"math/rand"
)

var vocab = map[string]int{"我": 0, "爱": 1, "自然": 2, "语言": 3, "处理": 4}
var idToWord = []string{"我", "爱", "自然", "语言", "处理"}

const vocabSize, embedDim = 5, 3

var (
	centerEmbed  = [vocabSize][embedDim]float64{}
	contextEmbed = [vocabSize][embedDim]float64{}
)

func dot(a, b [embedDim]float64) float64 {
	var s float64
	for i := 0; i < embedDim; i++ {
		s += a[i] * b[i]
	}
	return s
}

// sigmoid(x) = 1/(1+e^{-x})
func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

// 负采样：给定中心词 ID 和正样本上下文 ID，随机采 k 个负样本 ID（不等于正样本）
func negativeSample(centerID, positiveID, k int) []int {
	neg := make([]int, 0, k)
	for len(neg) < k {
		id := rand.Intn(vocabSize)
		if id != positiveID && id != centerID {
			neg = append(neg, id)
		}
	}
	return neg
}

func main() {
	rand.Seed(42)
	centerID, positiveID := 1, 2 // 中心词「爱」，正样本「自然」
	negIDs := negativeSample(centerID, positiveID, 3)

	vc := centerEmbed[centerID]
	posScore := sigmoid(dot(vc, contextEmbed[positiveID]))

	fmt.Printf("正样本 (爱, 自然) 得分 sigmoid=%.4f\n", posScore)
	fmt.Printf("负样本 ID: %v\n", negIDs)
	for _, id := range negIDs {
		s := sigmoid(dot(vc, contextEmbed[id]))
		fmt.Printf("  (爱, %s) sigmoid=%.4f\n", idToWord[id], s)
	}
	// 训练目标：正样本 sigmoid 接近 1，负样本接近 0
}
```

---

## 小结

| 步骤 | 含义 |
|------|------|
| 词表 + ID | 把词映射到 0..V-1 |
| 嵌入矩阵 | 词 ID → 稠密向量，核心可学习参数 |
| Skip-gram | 用中心词向量与上下文词向量内积 + softmax 做多分类 |
| 负采样 | 用正样本 + 随机负样本做二分类，近似多分类，便于大规模词表 |

实际实现还会涉及：滑动窗口取 (中心词, 上下文词)、学习率与迭代、词频与负采样分布等。上面几段 Go 代码仅用于理解「词 → 向量 → 预测/得分」这一主线原理。
