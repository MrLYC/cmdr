package core

import (
	"context"
	"errors"

	eventbus "github.com/asaskevich/EventBus"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/index"
	"github.com/asdine/storm/v3/q"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"logur.dev/logur"
)

type testCommandManager struct{}

func (testCommandManager) Close() error { return nil }
func (testCommandManager) Provider() CommandProvider {
	return CommandProviderDefault
}
func (testCommandManager) Query() (CommandQuery, error) { return nil, nil }
func (testCommandManager) Define(string, string, string) (Command, error) {
	return nil, nil
}
func (testCommandManager) Undefine(string, string) error { return nil }
func (testCommandManager) Activate(string, string) error { return nil }
func (testCommandManager) Deactivate(string) error       { return nil }

type testInitializer struct{}

func (testInitializer) Init(bool) error { return nil }

type testCmdrSearcher struct{}

func (testCmdrSearcher) GetReleaseAsset(context.Context, string, string) (CmdrReleaseAsset, error) {
	return CmdrReleaseAsset{Name: "cmdr"}, nil
}

type testDatabase struct{}

func (testDatabase) Close() error { return nil }
func (testDatabase) Init(interface{}) error {
	return nil
}
func (testDatabase) Save(interface{}) error {
	return nil
}
func (testDatabase) Update(interface{}) error {
	return nil
}
func (testDatabase) UpdateField(interface{}, string, interface{}) error {
	return nil
}
func (testDatabase) DeleteStruct(interface{}) error {
	return nil
}
func (testDatabase) Drop(interface{}) error {
	return nil
}
func (testDatabase) ReIndex(interface{}) error {
	return nil
}
func (testDatabase) One(string, interface{}, interface{}) error {
	return nil
}
func (testDatabase) All(interface{}, ...func(*index.Options)) error {
	return nil
}
func (testDatabase) AllByIndex(string, interface{}, ...func(*index.Options)) error {
	return nil
}
func (testDatabase) Find(string, interface{}, interface{}, ...func(q *index.Options)) error {
	return nil
}
func (testDatabase) Select(...q.Matcher) storm.Query {
	return nil
}
func (testDatabase) Count(interface{}) (int, error) {
	return 0, nil
}
func (testDatabase) Prefix(string, string, interface{}, ...func(*index.Options)) error {
	return nil
}
func (testDatabase) Range(string, interface{}, interface{}, interface{}, ...func(*index.Options)) error {
	return nil
}

var _ = Describe("Core", func() {
	It("should get and set configuration", func() {
		previous := GetConfiguration()
		defer SetConfiguration(previous)

		cfg := NewConfiguration()
		SetConfiguration(cfg)
		Expect(GetConfiguration()).To(Equal(cfg))
		Expect(cfg).To(BeAssignableToTypeOf(&viper.Viper{}))
	})

	It("should create exit errors", func() {
		err := NewExitError("failed", 42)
		Expect(err.Code()).To(Equal(42))
		Expect(err.Error()).To(Equal("failed: 42"))
	})

	It("should register command manager factories", func() {
		previous := factoriesCommandManager
		defer func() { factoriesCommandManager = previous }()
		factoriesCommandManager = make(map[CommandProvider]func(Configuration) (CommandManager, error))

		Expect(GetCommandManagerFactory(CommandProviderDefault)).To(BeNil())
		_, err := NewCommandManager(CommandProviderDefault, viper.New())
		Expect(err).To(Equal(ErrCommandManagerFactoryeNotFound))

		RegisterCommandManagerFactory(CommandProviderDefault, func(Configuration) (CommandManager, error) {
			return testCommandManager{}, nil
		})
		Expect(GetCommandManagerFactory(CommandProviderDefault)).NotTo(BeNil())
		manager, err := NewCommandManager(CommandProviderDefault, viper.New())
		Expect(err).NotTo(HaveOccurred())
		Expect(manager.Provider()).To(Equal(CommandProviderDefault))
	})

	It("should register initializer factories", func() {
		previous := factoriesInitializer
		defer func() { factoriesInitializer = previous }()
		factoriesInitializer = make(map[string]factoryInitializer)

		_, err := NewInitializer("missing", viper.New())
		Expect(err).To(Equal(ErrInitializerFactoryeNotFound))

		RegisterInitializerFactory("test", func(Configuration) (Initializer, error) {
			return testInitializer{}, nil
		})
		initializer, err := NewInitializer("test", viper.New())
		Expect(err).NotTo(HaveOccurred())
		Expect(initializer.Init(false)).To(Succeed())
	})

	It("should register cmdr searcher factories", func() {
		previous := factoriesCmdrSearcher
		defer func() { factoriesCmdrSearcher = previous }()
		factoriesCmdrSearcher = make(map[CmdrSearcherProvider]factoryCmdrSearcher)

		_, err := NewCmdrSearcher(CmdrSearcherProviderDefault, viper.New())
		Expect(err).To(Equal(ErrCmdrSearcherFactoryeNotFound))

		RegisterCmdrSearcherFactory(CmdrSearcherProviderDefault, func(Configuration) (CmdrSearcher, error) {
			return testCmdrSearcher{}, nil
		})
		searcher, err := NewCmdrSearcher(CmdrSearcherProviderDefault, viper.New())
		Expect(err).NotTo(HaveOccurred())
		asset, err := searcher.GetReleaseAsset(context.Background(), "latest", "asset")
		Expect(err).NotTo(HaveOccurred())
		Expect(asset.Name).To(Equal("cmdr"))
	})

	It("should register database models and factories", func() {
		previousModels := databaseModels
		previousFactory := databaseFactory
		defer func() {
			databaseModels = previousModels
			databaseFactory = previousFactory
		}()
		databaseModels = make(map[ModelType]interface{})
		databaseFactory = func() (Database, error) {
			return nil, ErrDatabaseFactoryNotSet
		}

		RegisterDatabaseModel(ModelTypeCommand, "command")
		Expect(GetDatabaseModel(ModelTypeCommand)).To(Equal("command"))
		Expect(GetDatabaseModels()).To(HaveKey(ModelTypeCommand))

		db, err := GetDatabase()
		Expect(db).To(BeNil())
		Expect(err).To(Equal(ErrDatabaseFactoryNotSet))

		SetDatabaseFactory(func() (Database, error) {
			return testDatabase{}, nil
		})
		Expect(GetDatabaseFactory()).NotTo(BeNil())
		db, err = GetDatabase()
		Expect(err).NotTo(HaveOccurred())
		Expect(db.Close()).To(Succeed())
	})

	It("should publish and subscribe events", func() {
		previous := eventBus
		defer func() { eventBus = previous }()
		eventBus = eventbus.New()

		count := 0
		SubscribeEvent("topic", func(delta int) {
			count += delta
		})
		SubscribeEventOnce("topic", func(delta int) {
			count += delta * 10
		})

		PublishEvent("topic", 1)
		PublishEvent("topic", 1)
		Expect(count).To(Equal(12))
	})

	It("should panic on invalid event subscribers", func() {
		previous := eventBus
		defer func() { eventBus = previous }()
		eventBus = eventbus.New()

		Expect(func() { SubscribeEvent("topic", "not-a-function") }).To(Panic())
		Expect(func() { SubscribeEventOnce("topic", "not-a-function") }).To(Panic())
	})

	It("should manage loggers and format fields", func() {
		previous := GetLogger()
		defer SetLogger(previous)

		SetLogger(logur.NoopLogger{})
		Expect(GetLogger()).To(Equal(logur.NoopLogger{}))

		logger := &terminalLogger{level: logur.Trace, errorKey: "error"}
		Expect(logger.getFieldsMessages(nil)).To(BeEmpty())
		Expect(logger.getFieldsMessages([]map[string]interface{}{{"name": "cmdr", "error": errors.New("boom")}})).To(ConsistOf("name=cmdr", "error=boom"))
		logger.withErrorStack = true
		Expect(logger.getFieldsMessages([]map[string]interface{}{{"error": errors.New("boom")}})[0]).To(ContainSubstring("boom"))

		logger.Trace("trace")
		logger.Debug("debug")
		logger.Info("info")
		logger.Warn("warn")
		logger.Error("error")
		logger.Info("")

		InitTerminalLogger(logur.Info, true, "error")
		Expect(GetLogger()).NotTo(BeNil())
	})
})
