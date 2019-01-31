PROJECTS = cpf2svl
CLEAN = $(addsuffix _clean, $(PROJECTS))

.PHONY: all $(PROJECTS)
all: $(PROJECTS)

$(PROJECTS):
	$(MAKE) -C $@


.PHONY: clean

clean: $(CLEAN)

%_clean:
	$(MAKE) -C $* clean
