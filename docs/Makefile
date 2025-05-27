# Minimal makefile for Sphinx documentation
#

# You can set these variables from the command line, and also
# from the environment for the first two.
SPHINXOPTS    ?=
SPHINXBUILD   ?= sphinx-build
SOURCEDIR     = .
BUILDDIR      = _build

# Put it first so that "make" without argument is like "make help".
help:
	@$(SPHINXBUILD) -M help "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)

.PHONY: help Makefile

# Catch-all target: route all unknown targets to Sphinx using the new
# "make mode" option.  $(O) is meant as a shortcut for $(SPHINXOPTS).
%: Makefile
	@$(SPHINXBUILD) -M $@ "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)


html: apidoc
	@rm -rf $(SPHINXBUILD)
	@$(SPHINXBUILD) -M $@ "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)
	for i in $(wildcard $(BUILDDIR)/html/_api/**.html); do \
	sed \
	-e 's@<span class="sig-paren">)</span>@<span class="sig-paren">]</span>@g' \
	-e 's@<span class="sig-paren">(</span>@<span class="sig-paren"> [</span>@g' \
	$$i > temp; \
	mv temp $$i; \
	done
apidoc:
	rm -rf _api
	rm -rf FMOe
	cp -r ../fmoe FMOe
	find FMOe -name "*.svl" -exec python svl2py.py {} \;
	sphinx-apidoc -f -o _api -M "FMOe"

deploy: 
	make clean
	git checkout gh-pages
	git reset --hard develop
	make html
	git add _build/html
	git commit -m "add html"
	git reset --hard $$(cd .. && git subtree split --prefix docs/_build/html/ gh-pages)
	git push -f origin gh-pages
	git checkout develop
