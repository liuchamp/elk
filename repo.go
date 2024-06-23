package elk

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/liuchamp/elk/outer/bo"
	"github.com/liuchamp/elk/outer/def"
	"github.com/liuchamp/elk/outer/dto"
	"github.com/liuchamp/elk/outer/imp"
	"github.com/liuchamp/elk/outer/router"
	"path/filepath"
)

type (
	// RepoConfig repo 输出配置
	RepoConfig struct {
		out    string // 输出路径， 相对位置
		dtoPre string // dto.go 前缀
		voPre  string // vo 前缀

		cache string // 缓存位置，只在运行中生成
	}
)

func newRepoConfig() RepoConfig {
	return RepoConfig{out: "repo"}
}

func RepoGenerator(c RepoConfig) gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			// Let ent create all the files.
			if err := next.Generate(g); err != nil {
				return err
			}
			c.cache = filepath.Join(g.Config.Target, "repo")

			//// 2. 输出dto
			err := dto.DtoOuter(g, c.cache)
			if err != nil {
				return err
			}
			//// 3. 输出 bo
			err = bo.BoOuter(g, c.cache)
			if err != nil {
				return err
			}
			//// 4. 输出 def
			err = def.DefOuter(g, c.cache)
			if err != nil {
				return err
			}
			//// 5. 输出 dao 层的imp
			err = imp.ImpOuter(g, c.cache)
			if err != nil {
				return err
			}
			// 6. 输出 router 层的imp
			err = router.RouterGen(g, c.cache)
			if err != nil {
				fmt.Printf("cat not gen router %s", err.Error())
				return err
			}
			return nil
		})
	}
}
