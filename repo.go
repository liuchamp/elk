package elk

import (
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk/outer/bo"
	"github.com/masseelch/elk/outer/def"
	"github.com/masseelch/elk/outer/dto"
	"github.com/masseelch/elk/outer/imp"
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

			// 2. 输出dto
			err := dto.DtoOuter(g, c.cache)
			if err != nil {
				return err
			}
			// 3. 输出 bo
			err = bo.BoOuter(g, c.cache)
			if err != nil {
				return err
			}
			// 4. 输出 def
			err = def.DefOuter(g, c.cache)
			if err != nil {
				return err
			}
			// 5. 输出 dao 层的imp
			err = imp.ImpOuter(g, c.cache)
			if err != nil {
				return err
			}
			return nil
		})
	}
}