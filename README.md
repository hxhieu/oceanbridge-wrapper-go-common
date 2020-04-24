# Shared modules

- Use **[wsdl2go](https://github.com/fiorix/wsdl2go)** command-line to generate the wrappers
- Then update the codegen `.RoundTripSoap12()` to `.RoundTripWithAction()`, it seems OceanBridge API not happy with 1.2 Content-Type generate by **wsdl2go**
