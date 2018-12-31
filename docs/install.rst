Installation
============

Via go standard tools
---------------------
As netorcai is implemented in Go_, it can be built with the `go command`_.
Installation steps are as follows.

1. Install a recent Go_ version.
2. Run :code:`go get github.com/netorcai/netorcai/cmd/netorcai`.
   This will download and compile netorcai.
   The executable will be put into ``${GOPATH}/bin``
   (if the ``GOPATH`` environment variable is unset,
   it should default to ``${HOME}/go`` or ``%USERPROFILE%\go``).

In brief.

.. code:: bash

    go get github.com/netorcai/netorcai/cmd/netorcai
    ${GOPATH:-${HOME}/go}/bin/netorcai --help

Via Nix
-------
Nix_ is a package manager with amazing properties that is available on
Linux-like systems.
It stores all the packages in a dedicated directory (usually :code:`/nix/store`),
which avoids interfering with classical system packages (usually in :code:`/usr`).

Once Nix is installed on your machine (instructions on `Nix's web page <Nix_>`_),
packages can be installed with :code:`nix-env --install` (:code:`-i`).
The following command shows how to install netorcai with Nix.

.. code:: bash

    # Install latest release.
    nix-env -f https://github.com/netorcai/netorcaipkgs/archive/master.tar.gz -iA netorcai

    # Alternatively, install latest commit.
    nix-env -f https://github.com/netorcai/netorcaipkgs/archive/master.tar.gz -iA netorcai_dev

.. _Go: https://golang.org/
.. _go command: https://golang.org/cmd/go/
.. _Nix: https://nixos.org/nix/
