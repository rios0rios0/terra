package repositorybuilders

import testkit "github.com/rios0rios0/testkit/pkg/test"

// cloneBase deep-copies base. [testkit.BaseBuilder.Clone] hands its copy back
// as a [testkit.Builder], which is always a *testkit.BaseBuilder.
func cloneBase(base *testkit.BaseBuilder) *testkit.BaseBuilder {
	clone, ok := base.Clone().(*testkit.BaseBuilder)
	if !ok {
		panic("testkit.BaseBuilder.Clone returned another kind of builder")
	}
	return clone
}
