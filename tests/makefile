MOE := moe.exe
MOEBATCH := moebatch.exe

.PHONY: test
test:
	rm -rf temp ||:
	mkdir temp
	python test_fragmentation.py

.PHONY: test_view
test_view:
	$(MOE) -exec "run ['../fmoe/presenter/fragmentation.svl', [test: 'fragmentation_gui', moe: 'resources\\\\\\test.moe']]"

.PHONY: test_view_merge
test_view_merge:
	$(MOE) -exec "run ['../fmoe/presenter/fragmentation.svl', [test: 'fragmentation_gui', moe: 'resources\\\\\\test.moe', merge: 1, merge_list: [[1,2], [3,4,5]]]]"

.PHONY: test_view_autofrag
test_view_autofrag:
	$(MOE) -exec "run ['../fmoe/presenter/fragmentation.svl', [test: 'fragmentation_gui', moe: 'autofrag.pdb', autofrag: 1, autofrag_svlb: 'autofrag.svlb']]"

.PHONY: test_view_hybrid
test_view_hybrid:
	$(MOE) -exec "run ['../fmoe/presenter/fragmentation.svl', [test: 'hybrid', moe: 'test_hybrid.pdb', hybrid: 1, autofrag: 'test_hybrid.autofrag']]"

.PHONY: test_merge
test_merge:
	$(MOE) -exec "run ['../fmoe/presenter/fragmentation.svl', [test: 'merge', merge_list: []]]"

.PHONY: clean
clean:
	rm temp/* ||:
