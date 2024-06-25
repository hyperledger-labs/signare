<h1 align="center">Hyperledger Signare</h1>
<p align="center"><i>An enterprise grade digital signing solution for DLT related applications and Ethereum clients.</i></p>
<div align="center">
<a href="https://github.com/hyperledger-labs/signare/pulls"><img src="https://img.shields.io/github/issues-pr/hyperledger-labs/signare" alt="Pull Requests Badge"/></a>
<a href="https://github.com/hyperledger-labs/signare/issues"><img src="https://img.shields.io/github/issues/hyperledger-labs/signare" alt="Issues Badge"/></a>
<a href="https://github.com/hyperledger-labs/signare/LICENSE"><img src="https://img.shields.io/github/license/hyperledger-labs/signare?color=2b9348" alt="License Badge"/></a>
</div>


Signare, a Hyperledger Lab, is an enterprise grade digital signing solution for DLT related applications and Ethereum clients. The application provides a REST API server to manage resource configuration and an ETH-JSON-RPC 2.0 server that provides functionality for generating, removing, listing and signing Ethereum transactions.

## 🗒 Contents
- [Scope of Lab](#scope-of-lab)
- [Useful Links](#globe_with_meridians-useful-links)
- [Features](#star-features)
- [Installation](#wrench-installation)
- [Getting Started](#astronaut-getting-started)
- [Contribute](#woman_technologist-contribute)
- [License](#pencil-license)
- [Initial committers](#initial-committers)
- [Sponsor](#sponsor)

## Scope of Lab

A security concern shared by most users of DLT applications is "keeping their private key private". In the enterprise space FIPS 140 is often used to inform institutions of how they must manage their private keys. Specifically, FIPS 140-2 Level 2 adds requirements for physical tamper-evidence (and/or tamper-resistance) and role-based authentication, which necessitates the use of an HSM or Cloud HSM.

The purpose of Hyperledger Signare is to provide a FIPS 140-2 Level 2 compliant signing solution for enterprise applications where various HSM and Cloud HSM vendors will be supported via plugins. Hyperledger Signare also provides role-based access controlled interfaces to solve multiple usecases, such as signing Ethereum transactions and for blockchain clients such as Hyperledger Besu to store keys in an HSM or Cloud HSM.

## :globe_with_meridians: Useful Links

- [Documentation](app/docs/mkdocs/docs/index.md): Discover the signare functionality and learn to configure it properly.
- [Changelog](app/docs/mkdocs/docs/CHANGELOG.md): Take a look at the record of changes.
- Feedback: Your help is key to develop the signare.
    - Found a bug? Need help fixing a problem? You can submit your issues [here](https://github.com/hyperledger-labs/signare/issues).

## :star: Features

signare comes with a range of features tailored for web3 integration purposes:

- **Different HSMs support**: signare support different types of HSM.
- **Ethereum accounts management**: generate new accounts, store and assign them to users.
- **Ethereum's transaction signing**: sign Ethereum transactions on the fly using managed accounts.

## :wrench: Installation

To start working with the signare, follow these installation steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/hyperledger-labs/signare.git
   ```

2. Navigate to the project's ``deployment`` directory:

   ```bash
   cd deployment
   ```

3. Build the binary:

   ```bash
   make build
   ```

4. Change directory to ``bin`` where the built binary will be stored:

   ```bash
   cd bin  
   ```

## :astronaut: Getting Started

To get started quickly, please refer to the [Getting started guide](app/docs/mkdocs/docs/getting-started/getting-started.md).
It covers a fast introduction to the signare's API usage, from creating a user with a new account to signing a transaction.

## :woman_technologist: Contribute

Contributions are always welcome! Please check our [Contribution guidelines](app/docs/mkdocs/docs/contribute/index.md) for details on how to get involved in the project's development.

## :pencil: License

This project is licensed under the [![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) license.

## Initial Committers

https://github.com/nano-adhara  
https://github.com/chookly314  
https://github.com/Jserrano27  
https://github.com/gynura  
https://github.com/ArturoGarciaRegueiro  
https://github.com/mkrielza  
https://github.com/coeniebeyers  

## Sponsor

- Susumu Toriumi (susumu.toriumi@datachain.jp) - Maintainer, YUI
