// Add current year to footer
document.addEventListener('DOMContentLoaded', function() {
    const currentYear = new Date().getFullYear();
    const footerYear = document.querySelector('footer p');
    if (footerYear) {
        footerYear.innerHTML = footerYear.innerHTML.replace('{{ current_year }}', currentYear);
    }

    // Add syntax highlighting for code blocks (can be enhanced in later phases)
    const codeBlocks = document.querySelectorAll('pre code');
    if (codeBlocks.length > 0) {
        // Simple syntax highlighting for Go code
        codeBlocks.forEach(function(block) {
            const code = block.innerHTML;

            // Highlight keywords
            const keywords = ['package', 'import', 'func', 'var', 'const', 'type', 'struct', 'interface',
                'if', 'else', 'for', 'range', 'switch', 'case', 'default', 'return', 'break',
                'continue', 'goto', 'map', 'chan', 'go', 'select', 'defer'];

            let highlighted = code;

            // Simple replacement for keywords (a real syntax highlighter would be better)
            // keywords.forEach(function(keyword) {
            //     const regex = new RegExp('\\b' + keyword + '\\b', 'g');
            //     highlighted = highlighted.replace(regex, '<span class="keyword">' + keyword + '</span>');
            // });

            // Simple string highlighting
            highlighted = highlighted.replace(/"[^"]*"/g, function(match) {
                return '<span class="string">' + match + '</span>';
            });

            // Simple comment highlighting
            highlighted = highlighted.replace(/\/\/[^\n]*/g, function(match) {
                return '<span class="comment">' + match + '</span>';
            });

            block.innerHTML = highlighted;
        });

        // Add styles for syntax highlighting
        const style = document.createElement('style');
        style.textContent = `
            pre code .keyword { color: #0076c0; }
            pre code .string { color: #42b983; }
            pre code .comment { color: #999; }
        `;
        document.head.appendChild(style);
    }
});