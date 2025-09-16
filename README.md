
```
                    .__            .___            .__                 __    
                    |__| ____    __| _/____ ___  __|  |   ____ _____  |  | __
                    |  |/    \  / __ |/ __ \\  \/  /  | _/ __ \\__  \ |  |/ /
                    |  |   |  \/ /_/ \  ___/ >    <|  |_\  ___/ / __ \|    < 
                    |__|___|  /\____ |\___  >__/\_ \____/\___  >____  /__|_ \
                            \/      \/    \/      \/         \/     \/     \/
```

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![MCP Protocol](https://img.shields.io/badge/MCP-Compatible-4A90E2?style=for-the-badge)](https://modelcontextprotocol.io/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Security](https://img.shields.io/badge/Security-OSINT-red?style=for-the-badge)](https://github.com/riza/indexleak-scanner)

</div>

> **A powerful MCP (Model Context Protocol) server for discovering and analyzing exposed directory listings on web servers. Perfect for security researchers, penetration testers, and cybersecurity professionals.**

## 🚀 Features

- **🔍 Directory Enumeration**: Automatically discover and scan exposed directory listings
- **📊 Risk Assessment**: Intelligent classification and risk scoring of discovered files
- **🎯 Pagination Support**: Handle large directories with efficient pagination
- **⚡ Real-time Analysis**: Live scanning and classification of sensitive data
- **🛡️ Security Focus**: Built specifically for cybersecurity and OSINT operations

## 🎯 Use Cases

- **Security Audits**: Identify exposed sensitive files and directories
- **Penetration Testing**: Discover potential attack vectors through directory traversal
- **OSINT Research**: Gather intelligence from publicly accessible directories
- **Compliance Checks**: Verify that sensitive data isn't publicly exposed
- **Vulnerability Assessment**: Systematic analysis of web server misconfigurations

## 🛠️ Installation

### Prerequisites

- Go 1.23 or higher
- Cursor IDE (for MCP integration)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/riza/indexleak-scanner.git
cd indexleak-scanner

# Build the binary
go build -o indexleak-scanner main.go

# Make it executable
chmod +x indexleak-scanner
```

## 📋 Configuration

### Setting up with Cursor IDE

1. **Add to MCP Configuration**

   Edit your Cursor MCP configuration file (`~/.cursor/mcp.json`):

   ```json
   {
     "mcpServers": {
       "indexleak": {
         "command": "/path/to/your/indexleak-scanner/indexleak-scanner",
         "args": [""],
         "env": {}
       }
     }
   }
   ```

2. **Restart Cursor IDE**

   After adding the configuration, restart Cursor to load the MCP server.

3. **Verify Installation**

   You should now see the IndexLeak tools available in your Cursor IDE's MCP tools panel.

## 🎮 Usage

### Available Tools

#### 1. `enter_directory`
Scans and analyzes directory listings from web servers.

**Parameters:**
- `url` (required): Target URL to scan
- `page` (optional): Page number for pagination (default: 1)
- `page_size` (optional): Items per page (default: 30)

**Example:**
```
URL: http://example.com/files/
Page: 1
Page Size: 50
```

## 🎯 Advanced Security Analysis

This tool can be used with detailed instructions in the `prompts/prompt-large.md` file for comprehensive security analysis. This file includes:

- **Systematic Scanning Protocol**: Step-by-step guide ensuring 100% coverage of all directories
- **Risk Classification Matrix**: Detailed risk scoring system based on file types
- **Security Assessment Template**: Professional security reporting format
- **Compliance Analysis**: Assessment from GDPR and other data protection regulations perspective

For detailed security analysis, you can use it by replacing the **TARGET HERE** section in the `prompts/prompt-large.md` file with your target URL.

### 🔍 Example Workflow in Cursor

1. **Start a Security Scan**
   ```
   Use the enter_directory tool to scan: http://target-server.com/
   ```

2. **Analyze Results**
   The tool will automatically:
   - Classify each file and directory
   - Assign risk scores (1-10)
   - Identify sensitive data exposure


## 📊 Risk Classification

The tool uses an intelligent risk scoring system:

| Score | Level | Description |
|-------|-------|-------------|
| **10** | 🔴 **Critical** | Cryptographic keys, database files, credentials |
| **8-9** | 🟠 **High** | Configuration files, sensitive documents |
| **6-7** | 🟡 **Medium** | Business documents, archives |
| **4-5** | 🔵 **Low** | Media files, executables |
| **1-3** | 🟢 **Minimal** | Public content, documentation |

## 📈 Sample Output

```
Directory contents (http://example.com/files/):

Page 1 of 3 (Total entries: 87, Entries per page: 30)

FILE: config.php (http://example.com/files/config.php)
FILE: database_backup.sql (http://example.com/files/database_backup.sql)
DIRECTORY: admin (http://example.com/files/admin/)
FILE: users.csv (http://example.com/files/users.csv)
...
```

## 🔒 Security Considerations

- **Ethical Use Only**: This tool is designed for legitimate security testing and research
- **Permission Required**: Always ensure you have proper authorization before scanning
- **Responsible Disclosure**: Report findings through appropriate channels
- **Legal Compliance**: Follow applicable laws and regulations in your jurisdiction

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

This tool is intended for educational and authorized security testing purposes only. Users are responsible for ensuring they have proper permission before scanning any systems. The authors are not responsible for any misuse or damage caused by this tool.

## 🙏 Acknowledgments

- Built with the [Model Context Protocol (MCP)](https://modelcontextprotocol.io/)
- Powered by [Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk)
- Designed for integration with [Cursor IDE](https://cursor.sh/)

## 📞 Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/riza/indexleak-scanner/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/riza/indexleak-scanner/discussions)
- 📧 **Contact**: [rizasabuncu[at]outlook.com](mailto:rizasabuncu@outlook.com)

---

<div align="center">

Made with ❤️ for the cybersecurity community

</div>
