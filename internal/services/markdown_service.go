package services

import (
	"bufio"
	"context"
	"nexzap/templates/partials"
	"regexp"
	"strings"
)

// MarkdownParser holds the state and configuration for parsing Markdown
type MarkdownParser struct {
	output         strings.Builder
	paragraphLines []string
	codeLines      []string
	inCodeBlock    bool
	language       string
}

// NewMarkdownParser creates a new MarkdownParser instance
func NewMarkdownParser() *MarkdownParser {
	return &MarkdownParser{}
}

// InlinePatterns holds regular expressions for inline Markdown elements
type InlinePatterns struct {
	code   *regexp.Regexp
	bold   *regexp.Regexp
	italic *regexp.Regexp
	link   *regexp.Regexp
}

// NewInlinePatterns initializes regular expressions for inline parsing
func NewInlinePatterns() *InlinePatterns {
	return &InlinePatterns{
		code:   regexp.MustCompile("`(.+?)`"),
		bold:   regexp.MustCompile(`\*\*(.+?)\*\*`),
		italic: regexp.MustCompile(`\*(.+?)\*`),
		link:   regexp.MustCompile(`\[(.+?)\]\((.+?)\)`),
	}
}

// parseInline processes inline Markdown elements and applies Tailwind classes
func (p *MarkdownParser) parseInline(text string, patterns *InlinePatterns) string {
	result := text

	// Process code spans
	result = patterns.code.ReplaceAllStringFunc(result, func(match string) string {
		code := patterns.code.FindStringSubmatch(match)[1]
		return `<code class="bg-gray-100 p-1 rounded">` + code + `</code>`
	})

	// Process bold text
	result = patterns.bold.ReplaceAllStringFunc(result, func(match string) string {
		bold := patterns.bold.FindStringSubmatch(match)[1]
		var res strings.Builder
		partials.Bold(bold).Render(context.Background(), &res)
		return res.String()
	})

	// Process italic text
	result = patterns.italic.ReplaceAllStringFunc(result, func(match string) string {
		italic := patterns.italic.FindStringSubmatch(match)[1]
		var res strings.Builder
		partials.Italic(italic).Render(context.Background(), &res)
		return res.String()
	})

	// Process links
	result = patterns.link.ReplaceAllStringFunc(result, func(match string) string {
		parts := patterns.link.FindStringSubmatch(match)
		var res strings.Builder
		partials.Link(parts[1], parts[2]).Render(context.Background(), &res)
		return res.String()
	})

	return result
}

// buildParagraphText constructs paragraph text, handling line breaks within paragraphs
func buildParagraphText(lines []string) string {
	var paragraph strings.Builder
	for i, line := range lines {
		trimmed := strings.TrimRight(line, " ")
		paragraph.WriteString(trimmed)
		if i < len(lines)-1 {
			paragraph.WriteString("<br>")
		}
	}
	return paragraph.String()
}

// flushParagraph writes accumulated paragraph lines to output
func (p *MarkdownParser) flushParagraph(patterns *InlinePatterns) {
	if len(p.paragraphLines) == 0 {
		return
	}
	paragraphText := buildParagraphText(p.paragraphLines)
	processed := p.parseInline(paragraphText, patterns)
	p.output.WriteString("<p>" + processed + "</p>\n")
	p.paragraphLines = nil
}

// processHeading handles Markdown heading lines
func (p *MarkdownParser) processHeading(line string, patterns *InlinePatterns) bool {
	if !strings.HasPrefix(line, "#") {
		return false
	}

	level := 0
	for i, char := range line {
		if char != '#' {
			if i < len(line) && line[i] == ' ' {
				level = i
				break
			}
			return false
		}
	}

	if level >= 1 && level <= 6 {
		text := line[level+1:]
		processed := p.parseInline(text, patterns)
		partials.Header(level, processed).Render(context.Background(), &p.output)
		return true
	}
	return false
}

// processCodeBlock handles the start/end of code blocks
func (p *MarkdownParser) processCodeBlock(line string) bool {
	if !strings.HasPrefix(line, "```") {
		return false
	}

	if !p.inCodeBlock {
		// Start code block
		p.inCodeBlock = true
		p.language = strings.TrimSpace(strings.TrimLeft(line, "`"))
	} else {
		// End code block
		p.inCodeBlock = false
		code := strings.Join(p.codeLines, "\n")
		// TODO : use language
		partials.Snippet(code).Render(context.Background(), &p.output)
		p.codeLines = nil
		p.language = ""
	}
	return true
}

// ParseMarkdown converts Markdown text to HTML with Tailwind CSS classes
func (p *MarkdownParser) ParseMarkdown(md string) string {
	scanner := bufio.NewScanner(strings.NewReader(md))
	patterns := NewInlinePatterns()

	for scanner.Scan() {
		line := scanner.Text()

		if p.inCodeBlock {
			if p.processCodeBlock(line) {
				continue
			}
			p.codeLines = append(p.codeLines, line)
			continue
		}

		if p.processCodeBlock(line) {
			continue
		}

		if p.processHeading(line, patterns) {
			p.flushParagraph(patterns)
			continue
		}

		if strings.TrimSpace(line) == "" {
			continue
		}

		p.paragraphLines = append(p.paragraphLines, line)
	}

	result := p.output.String()
	*p = *NewMarkdownParser()

	return result
}
