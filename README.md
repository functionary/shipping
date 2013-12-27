This package is licensed under the terms of the Apache Software License version 2.0. See accompanying file LICENSE-2.0.txt for the text of the license.

This package is not usable at this time without significant modification. It depends on another package for struct definitions, which I must move into this package so that it is fully functional.

The plan is to finish the FedEx (which is not yet included here), UPS, and USPS APIs and then implement an abstraction layer in the root directory, 'shipping.go'. This should enable users to shop and ship packages easily without having to talk to the individual APIs.
