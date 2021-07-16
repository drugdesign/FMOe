
.PHONY: test
test:
	python test_fragmentation.py

.PHONY: test_view
test_view:
	moe.exe -exec "run ['../fragmentation.svl', [test: 'fragmentation_gui', moe: 'resources\\\\\\test.moe']]"

.PHONY: test_view_merge
test_view_merge:
	moe.exe -exec "run ['../fragmentation.svl', [test: 'fragmentation_gui', moe: 'resources\\\\\\test.moe', merge: 1, merge_list: [[1,2], [3,4,5]]]]"

.PHONY: test_merge
test_merge:
	moe.exe -exec "run ['../fragmentation.svl', [test: 'merge', merge_list: []]]"

.PHONY: clean
clean:
	rm temp/* ||:
