package services_test

import (
	"nexzap/internal/services"
	"testing"
)

func TestParseMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single heading",
			input:    "# Heading 1",
			expected: "<h1>Heading 1</h1>",
		},
		{
			name:     "multiple headings",
			input:    "# H1\n## H2\n### H3",
			expected: "<h1>H1</h1><h2>H2</h2><h3>H3</h3>",
		},
		{
			name:     "paragraph with multiple lines",
			input:    "Line one\nLine two\nLine three",
			expected: "<p>Line one<br>Line two<br>Line three</p>",
		},
		{
			name:     "paragraph with inline elements",
			input:    "This is *italic*, **bold**, `code`, and [link](url).",
			expected: "<p>This is <em class=\"italic\">italic</em>, <strong class=\"font-bold\">bold</strong>, <code class=\"bg-gray-100 p-1 rounded\">code</code>, and <a href=\"url\" class=\"link link-primary\">link</a>.</p>",
		},
		{
			name:     "code block",
			input:    "```\nsome code\nhere\n```",
			expected: "<pre class=\"code bg-base-200 cm-s-daisyui\">some code\nhere</pre>",
		},
		{
			name:     "mixed content",
			input:    "Paragraph before.\n```\ncode block\n```\nParagraph after.",
			expected: "<p>Paragraph before.</p><pre class=\"code bg-base-200 cm-s-daisyui\">code block</pre><p>Paragraph after.</p>",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "only empty lines",
			input:    "\n\n\n",
			expected: "",
		},
		{
			name:     "paragraph with trailing spaces",
			input:    "Line with trailing spaces.  \nAnother line.",
			expected: "<p>Line with trailing spaces.<br>Another line.</p>",
		},
		{
			name: "complex content with headings, lists, and code block",
			input: "## Task: Calculate Sum with Goroutines and Channels\n" +
				"\n" +
				"### Instructions\n" +
				"\n" +
				"Write a function `CalculateSum(numbers []int) int` that will:\n" +
				"1. Take a slice of integers as input.\n" +
				"2. Split the work of summing the numbers into two goroutines if the slice has more than one element.\n" +
				"3. Use a channel to communicate partial sums from the goroutines and return the total sum.\n" +
				"\n" +
				"#### Steps:\n" +
				"- Declare a channel for partial sums.\n" +
				"- Split the slice and use goroutines for summing each part (if length > 1).\n" +
				"- Return the combined sum from the channel results.\n" +
				"\n" +
				"#### Note on Slices:\n" +
				"A slice in Go is a flexible view of an array. Use `len()` to get its length and `[:]` notation to create sub-slices (e.g., `numbers[:2]` for the first two elements, `numbers[2:]` for the rest).\n" +
				"\n" +
				"#### Example:\n" +
				"- `CalculateSum([]int{1, 2, 3, 4})` should return `10`.\n" +
				"- `CalculateSum([]int{5})` should return `5`.\n" +
				"- `CalculateSum([]int{})` should return `0`.",
			expected: "<h2>Task: Calculate Sum with Goroutines and Channels</h2>" +
				"<h3>Instructions</h3>" +
				"<p>Write a function <code class=\"bg-gray-100 p-1 rounded\">CalculateSum(numbers []int) int</code> that will:" +
				"<br>" +
				"1. Take a slice of integers as input." +
				"<br>" +
				"2. Split the work of summing the numbers into two goroutines if the slice has more than one element." +
				"<br>" +
				"3. Use a channel to communicate partial sums from the goroutines and return the total sum.</p>" +
				"<h4>Steps:</h4>" +
				"<p>- Declare a channel for partial sums." +
				"<br>- Split the slice and use goroutines for summing each part (if length > 1)." +
				"<br>- Return the combined sum from the channel results.</p>" +
				"<h4>Note on Slices:</h4>" +
				"<p>A slice in Go is a flexible view of an array. Use " +
				"<code class=\"bg-gray-100 p-1 rounded\">len()</code> to get its length and " +
				"<code class=\"bg-gray-100 p-1 rounded\">[:]</code> notation to create sub-slices (e.g., " +
				"<code class=\"bg-gray-100 p-1 rounded\">numbers[:2]</code> for the first two elements, " +
				"<code class=\"bg-gray-100 p-1 rounded\">numbers[2:]</code> for the rest).</p>" +
				"<h4>Example:</h4>" +
				"<p>- " +
				"<code class=\"bg-gray-100 p-1 rounded\">CalculateSum([]int{1, 2, 3, 4})</code> should return " +
				"<code class=\"bg-gray-100 p-1 rounded\">10</code>." +
				"<br>- " +
				"<code class=\"bg-gray-100 p-1 rounded\">CalculateSum([]int{5})</code> should return " +
				"<code class=\"bg-gray-100 p-1 rounded\">5</code>." +
				"<br>- " +
				"<code class=\"bg-gray-100 p-1 rounded\">CalculateSum([]int{})</code> should return " +
				"<code class=\"bg-gray-100 p-1 rounded\">0</code>.</p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := services.NewMarkdownParser()
			output := parser.ParseMarkdown(tt.input)
			if output != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, output)
			}
		})
	}
}
