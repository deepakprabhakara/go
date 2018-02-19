package litecoin

import (
	"testing"

	"github.com/ltcsuite/ltcd/chaincfg"
	"github.com/stretchr/testify/assert"
)

func TestAddressGenerator(t *testing.T) {
	// Generated using https://iancoleman.github.io/bip39/
	// Root key:
	// Ltpv71G8qDifUiNesHWSni4KA4G3kZBRBX9nwuqhw8Jt6ukdUwuYxsQcLbgKrWTE2FhmFGKeJgytys579WDexNEWTEVToNeHiKQV3PyeF69yiui
	// Derivation Path m/44'/2'/0'/0:
	// Ltpv78suy1WqwyejGuFfTkcm1emtsbHu69ksFKNuE5sHXyBJzTCARd8jH4Q7jqT8ozRv7oFEG53wruiZjeDQZpLbVMnyVx5cgdc3LnNfJ85RdfF
	// Ltub2a4FZnwPC8BgtmxViN796Ec7K9z6AQKXKqUvu4daU9xuapPEkYBAHfWNzEwpFjd2sHwX9kCFPxyMuWTHHiizViuaa27t85DP5ezobRwbYPm
	generator, err := NewAddressGenerator("Ltub2a4FZnwPC8BgtmxViN796Ec7K9z6AQKXKqUvu4daU9xuapPEkYBAHfWNzEwpFjd2sHwX9kCFPxyMuWTHHiizViuaa27t85DP5ezobRwbYPm", &chaincfg.MainNetParams)
	assert.NoError(t, err)

	expectedChildren := []struct {
		index   uint32
		address string
	}{
		{0, "Lhd98J63jWM44tY8tcGPcvCdRDruDadyJj"},
		{1, "LMy56RybHawcF2NKX9d2jnfuWd42Hn1AyX"},
		{2, "LU4ueCvSALxpDV35VFeY6XSwYrtTmrNC1S"},
		{3, "Led8FBxx4kMTUgsgcXQymafHiraoApXVmq"},
		{4, "LhNH4gjdyEXhNyah22FDwvaAeJwfxd7kng"},
		{5, "Lbe14PXkYNHhFLm6Npsu3oZboM6xDtjJv7"},
		{6, "LTgHff4zjZBoLcnPASc4mRW7ZureNzA2qM"},
		{7, "LeuBMeguAXEzDfqyxqR9pJnRUKaXmPdsRz"},
		{8, "LgyqCpKJQLvgLmXpSu918gniWyFHPoYgZ2"},
		{9, "LdPgyKjqzkNv8hRfkjSmNG1zbhUUogbBL9"},

		{100, "LNCv2h63oJf4qqgXnV4DiSJT5oWiKZMpV4"},
		{101, "LemDH6aWCFEewFeGYGYqRXF2KcpnEeUZAy"},
		{102, "LRQ4k126ZLTLBJ15Y1qEBzfdSa1uhRP6eS"},
		{103, "LQARQJFAG1zTEEovMmkdr4z8sG3w5TdyAA"},
		{104, "LdBkBT2PX4dmP1KfUVcKd2ypZLT7q8FRwU"},
		{105, "LhwWaQsC8ygGUsUr7tjADN8ywTqj7gwnFh"},
		{106, "LazYfFBTwbipsrMRY4xht49QwjzevF7Qz2"},
		{107, "Lh39GmoawktnPSFjd4ruUQzK154vDuFUXE"},
		{108, "LWm6UqQnXzgTKX3heA6A5k4fxSBc1SaG7T"},
		{109, "LKbGGCt4VS7hiCkCasYQAcCZB7Tzrq2LeY"},

		{1000, "LTwWRfBS2jUbf1VBR5kdVaxPgdYdMLzfNN"},
	}

	for _, child := range expectedChildren {
		address, err := generator.Generate(child.index)
		assert.NoError(t, err)
		assert.Equal(t, child.address, address)
	}
}
