package webgl

import (
	"errors"
	"github.com/tokkenno/seed/core"
	"github.com/tokkenno/seed/core/cameras"
	"github.com/tokkenno/seed/core/math"
	"github.com/tokkenno/seed/render"
	"github.com/tokkenno/seed/render/webgl/dom"
	"github.com/tokkenno/seed/render/webgl/gl"
)

type GeometryProgram struct {
	geometry
	program
	wireFrame bool
}

type Renderer struct {
	Canvas *dom.Canvas
	Context

	AutoClear        bool
	AutoClearColor   bool
	AutoClearDepth   bool
	AutoClearStencil bool

	SortObjects bool

	ClippingPlanes       []
	LocalClippingEnabled bool

	GammaFactor float64
	GammaInput  bool
	GammaOutput bool

	PhysicallyCorrentLights bool
	ToneMapping ToneMapping
	ToneMappingExposure   float64
	ToneMappingWhitePoint float64

	MaxMorphTargets int
	MaxMorphNormals int

	isContextLost bool

	frameBuffer

	currentRenderTarget
	currentFrameBuffer
	currentMaterialId int64

	currentGeometryProgram *GeometryProgram

	currentCamera *cameras.Camera
	currentArrayCamera

	currentViewport    *math.Vector4
	currentScissor     *math.Vector4
	currentScissorTest bool

	usedTextureUnits int

	width  int
	height int

	pixelRatio float32

	viewPort    *math.Vector4
	scissor     *math.Vector4
	scissorTest bool

	frustum              *math.Frustum
	clipping             *gl.Clipping
	clippingEnabled      bool
	localClippingEnabled bool

	projectScreenMatrix *math.Vector4

	vector3 *math.Vector3
}

func NewRenderer(options *RendererOptions) {
	renderer := new(Renderer)
	renderer.Canvas = options.Canvas
}

func (renderer *Renderer) Render(scene *core.Scene, camera *cameras.Camera, target *render.Target, forceClear bool) error {
	if camera == nil {
		return errors.New("the camera can't be null")
	}

	if renderer.isContextLost {
		return nil
	}

	// reset caching for this frame
	renderer.currentGeometryProgram.geometry = nil
	renderer.currentGeometryProgram.program = nil
	renderer.currentGeometryProgram.wireFrame = false
	renderer.currentMaterialId = -1
	renderer.currentCamera = nil

	// update scene graph
	if scene.AutoUpdateRender {
		scene.UpdateMatrixWorld(false)
	}

	// update camera matrices and frustum
	if camera.Parent == nil {
		camera.UpdateMatrixWorld(false)
	}

	renderer.currentRenderState = renderer.renderStates.Get(scene, camera)
	renderer.currentRenderState.Init()
}
