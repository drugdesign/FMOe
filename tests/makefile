
.PHONY: test
test:
	python test_fragmentation.py

.PHONY: test_view
test_view:
	moe.exe -exec "run ['../fragmentation.svl', [test: 'fragmentation_gui', moe: 'resources\\\\\\test.moe']]"

.PHONY: clean
clean:
	rm temp/* ||:
