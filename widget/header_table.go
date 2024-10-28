package widget

import (
	"math"
	"strconv"

	"github.com/djpken/go-fyne"
	"github.com/djpken/go-fyne/canvas"
	"github.com/djpken/go-fyne/driver/desktop"
	"github.com/djpken/go-fyne/internal/cache"
	"github.com/djpken/go-fyne/internal/widget"
	"github.com/djpken/go-fyne/theme"
)

// allTableCellsID represents all table cells whe refreshing requested cells

// Declare conformity with interfaces
var _ desktop.Cursorable = (*StaticTable)(nil)
var _ desktop.Hoverable = (*StaticTable)(nil)

// var _ fyne.Tappable = (*StaticTable)(nil)
var _ fyne.Widget = (*StaticTable)(nil)

type StaticTable struct {
	BaseWidget

	Length                                        func() (rows int, cols int)                      `json:"-"`
	CreateCell                                    func() fyne.CanvasObject                         `json:"-"`
	UpdateCell                                    func(id TableCellID, template fyne.CanvasObject) `json:"-"`
	OnSelected                                    func(id TableCellID)                             `json:"-"`
	OnUnselected                                  func(id TableCellID)                             `json:"-"`
	ShowHeaderRow                                 bool
	ShowHeaderColumn                              bool
	CreateHeader                                  func() fyne.CanvasObject                         `json:"-"`
	UpdateHeader                                  func(id TableCellID, template fyne.CanvasObject) `json:"-"`
	StickyRowCount                                int
	StickyColumnCount                             int
	cells                                         *headerTableCells
	columnWidths, rowHeights                      map[int]float32
	moveCallback                                  func()
	offset                                        fyne.Position
	content                                       *widget.Scroll
	cellSize, headerSize                          fyne.Size
	stuckXOff, stuckYOff, stuckWidth, stuckHeight float32
	top, left, corner, dividerLayer               *staticClip
	hoverHeaderRow, hoverHeaderCol                int
}

func NewStaticTable(length func() (rows int, cols int), create func() fyne.CanvasObject, update func(TableCellID, fyne.CanvasObject)) *StaticTable {
	t := &StaticTable{Length: length, CreateCell: create, UpdateCell: update}
	t.ShowHeaderColumn = true
	t.ShowHeaderRow = true
	t.ExtendBaseWidget(t)
	return t
}

// CreateRenderer returns a new renderer for the table.
//
// Implements: fyne.Widget
func (t *StaticTable) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)

	t.propertyLock.Lock()
	t.headerSize = t.createHeader().MinSize()
	if t.columnWidths != nil {
		if v, ok := t.columnWidths[-1]; ok {
			t.headerSize.Width = v
		}
	}
	if t.rowHeights != nil {
		if v, ok := t.rowHeights[-1]; ok {
			t.headerSize.Height = v
		}
	}
	t.cellSize = t.templateSize()
	t.cells = newHeaderTableCells(t)
	t.content = widget.NewScroll(t.cells)
	t.top = newHeaderClip(t, &fyne.Container{})
	t.left = newHeaderClip(t, &fyne.Container{})
	t.corner = newHeaderClip(t, &fyne.Container{})
	t.dividerLayer = newHeaderClip(t, &fyne.Container{})
	t.propertyLock.Unlock()

	r := &staticTableRender{t: t}
	r.SetObjects([]fyne.CanvasObject{t.top, t.left, t.corner, t.dividerLayer, t.content})
	t.content.OnScrolled = func(pos fyne.Position) {
		t.offset = pos
		impl := t.super()
		if impl == nil {
			return
		}
		t.propertyLock.Lock()
		t.themeCache = nil
		t.propertyLock.Unlock()
		render := cache.Renderer(t.cells)
		render.(*headerTableCellsRenderer).scroll()
	}

	r.Layout(t.Size())
	return r
}

func (t *StaticTable) Cursor() desktop.Cursor {
	if t.hoverHeaderRow != noCellMatch {
		return desktop.VResizeCursor
	} else if t.hoverHeaderCol != noCellMatch {
		return desktop.HResizeCursor
	}

	return desktop.DefaultCursor
}

func (t *StaticTable) MouseIn(ev *desktop.MouseEvent) {
	t.hoverAt(ev.Position)
}

func (t *StaticTable) MouseMoved(ev *desktop.MouseEvent) {
	t.hoverAt(ev.Position)
}

func (t *StaticTable) MouseOut() {
}

// RefreshItem refreshes a single item, specified by the item ID passed in.
//
// Since: 2.4
func (t *StaticTable) RefreshItem(id TableCellID) {
	if t.cells == nil {
		return
	}
	r := cache.Renderer(t.cells)
	if r == nil {
		return
	}

	r.(*headerTableCellsRenderer).refreshForID(id)
}

// SetColumnWidth supports changing the width of the specified column. Columns normally take the width of the template
// cell returned from the CreateCell callback. The width parameter uses the same units as a fyne.Size type and refers
// to the internal content width not including the divider size.
//
// Since: 1.4.1
func (t *StaticTable) SetColumnWidth(id int, width float32) {
	t.propertyLock.Lock()
	if id < 0 {
		if t.headerSize.Width == width {
			t.propertyLock.Unlock()
			return
		}
		t.headerSize.Width = width
	}

	if t.columnWidths == nil {
		t.columnWidths = make(map[int]float32)
	}

	if set, ok := t.columnWidths[id]; ok && set == width {
		t.propertyLock.Unlock()
		return
	}
	t.columnWidths[id] = width
	t.propertyLock.Unlock()

}

// SetRowHeight supports changing the height of the specified row. Rows normally take the height of the template
// cell returned from the CreateCell callback. The height parameter uses the same units as a fyne.Size type and refers
// to the internal content height not including the divider size.
//
// Since: 2.3
func (t *StaticTable) SetRowHeight(id int, height float32) {
	t.propertyLock.Lock()
	if id < 0 {
		if t.headerSize.Height == height {
			t.propertyLock.Unlock()
			return
		}
		t.headerSize.Height = height
	}

	if t.rowHeights == nil {
		t.rowHeights = make(map[int]float32)
	}

	if set, ok := t.rowHeights[id]; ok && set == height {
		t.propertyLock.Unlock()
		return
	}
	t.rowHeights[id] = height
	t.propertyLock.Unlock()
}

// columnAt returns a positive integer (or 0) for the column that is found at the `pos` X position.
// If the position is between cells the method will return a negative integer representing the next column,
// i.e. -1 means the gap between 0 and 1.
func (t *StaticTable) columnAt(pos fyne.Position) int {
	dataCols := 0
	if f := t.Length; f != nil {
		_, dataCols = t.Length()
	}

	visibleColWidths, offX, minCol, maxCol := t.visibleColumnWidths(t.cellSize.Width, dataCols)
	i := minCol
	end := maxCol
	if pos.X < t.stuckXOff+t.stuckWidth {
		offX = t.stuckXOff
		end = t.StickyColumnCount
		i = 0
	} else {
		pos.X += t.content.Offset.X
		offX += t.stuckXOff
	}
	padding := t.Theme().Size(theme.SizeNamePadding)
	for x := offX; i < end; x += visibleColWidths[i-1] + padding {
		if pos.X < x {
			return -i // the space between i-1 and i
		} else if pos.X < x+visibleColWidths[i] {
			return i
		}
		i++
	}
	return noCellMatch
}

func (t *StaticTable) createHeader() fyne.CanvasObject {
	if f := t.CreateHeader; f != nil {
		return f()
	}

	l := NewLabel("00")
	l.TextStyle.Bold = true
	l.Alignment = fyne.TextAlignCenter
	return l
}

func (t *StaticTable) hoverAt(pos fyne.Position) {
	col := t.columnAt(pos)
	row := t.rowAt(pos)
	overHeaderRow := t.ShowHeaderRow && pos.Y < t.headerSize.Height
	overHeaderCol := t.ShowHeaderColumn && pos.X < t.headerSize.Width
	if overHeaderRow && !overHeaderCol {
		if col >= 0 {
			t.hoverHeaderCol = noCellMatch
		} else {
			t.hoverHeaderCol = -col - 1
		}
	} else {
		t.hoverHeaderCol = noCellMatch
	}
	if overHeaderCol && !overHeaderRow {
		if row >= 0 {
			t.hoverHeaderRow = noCellMatch
		} else {
			t.hoverHeaderRow = -row - 1
		}
	} else {
		t.hoverHeaderRow = noCellMatch
	}

}

// rowAt returns a positive integer (or 0) for the row that is found at the `pos` Y position.
// If the position is between cells the method will return a negative integer representing the next row,
// i.e. -1 means the gap between rows 0 and 1.
func (t *StaticTable) rowAt(pos fyne.Position) int {
	dataRows := 0
	if f := t.Length; f != nil {
		dataRows, _ = t.Length()
	}

	visibleRowHeights, offY, minRow, maxRow := t.visibleRowHeights(t.cellSize.Height, dataRows)
	i := minRow
	end := maxRow
	if pos.Y < t.stuckYOff+t.stuckHeight {
		offY = t.stuckYOff
		end = t.StickyRowCount
		i = 0
	} else {
		pos.Y += t.content.Offset.Y
		offY += t.stuckYOff
	}
	padding := t.Theme().Size(theme.SizeNamePadding)
	for y := offY; i < end; y += visibleRowHeights[i-1] + padding {
		if pos.Y < y {
			return -i // the space between i-1 and i
		} else if pos.Y >= y && pos.Y < y+visibleRowHeights[i] {
			return i
		}
		i++
	}
	return noCellMatch
}

func (t *StaticTable) templateSize() fyne.Size {
	if f := t.CreateCell; f != nil {
		template := createItemAndApplyThemeScope(f, t) // don't use cache, we need new template
		if !t.ShowHeaderRow && !t.ShowHeaderColumn {
			return template.MinSize()
		}
		return template.MinSize().Max(t.createHeader().MinSize())
	}

	fyne.LogError("Missing CreateCell callback required for StaticTable", nil)
	return fyne.Size{}
}

func (t *StaticTable) updateHeader(id TableCellID, o fyne.CanvasObject) {
	if f := t.UpdateHeader; f != nil {
		f(id, o)
		return
	}

	l := o.(*Label)
	if id.Row < 0 {
		ids := []rune{'A' + rune(id.Col%26)}
		pre := (id.Col - id.Col%26) / 26
		for pre > 0 {
			ids = append([]rune{'A' - 1 + rune(pre%26)}, ids...)
			pre = (pre - pre%26) / 26
		}
		l.SetText(string(ids))
	} else if id.Col < 0 {
		l.SetText(strconv.Itoa(id.Row + 1))
	} else {
		l.SetText("")
	}
}

func (t *StaticTable) stickyColumnWidths(colWidth float32, cols int) (visible []float32) {
	if cols == 0 {
		return []float32{}
	}

	max := t.StickyColumnCount
	if max > cols {
		max = cols
	}

	visible = make([]float32, max)

	if len(t.columnWidths) == 0 {
		for i := 0; i < max; i++ {
			visible[i] = colWidth
		}
		return
	}

	for i := 0; i < max; i++ {
		height := colWidth

		if h, ok := t.columnWidths[i]; ok {
			height = h
		}

		visible[i] = height
	}
	return
}

func (t *StaticTable) visibleColumnWidths(colWidth float32, cols int) (visible map[int]float32, offX float32, minCol, maxCol int) {
	maxCol = cols
	colOffset, headWidth := float32(0), float32(0)
	visible = make(map[int]float32)

	if t.content.Size().Width <= 0 {
		return
	}

	padding := t.Theme().Size(theme.SizeNamePadding)
	stick := t.StickyColumnCount
	size := t.size.Load()

	if len(t.columnWidths) == 0 {
		paddedWidth := colWidth + padding

		offX = float32(math.Floor(float64(t.offset.X/paddedWidth))) * paddedWidth
		minCol = int(math.Floor(float64(offX / paddedWidth)))
		//maxCol = int(math.Ceil(float64((t.offset.X + size.Width) / paddedWidth)))

		if minCol > cols-1 {
			minCol = cols - 1
		}
		if minCol < 0 {
			minCol = 0
		}

		if maxCol > cols {
			maxCol = cols
		}

		visible = make(map[int]float32, maxCol-minCol+stick)
		for i := minCol; i < maxCol; i++ {
			visible[i] = colWidth
		}
		for i := 0; i < stick; i++ {
			visible[i] = colWidth
		}
		return
	}

	for i := 0; i < cols; i++ {
		width := colWidth
		if w, ok := t.columnWidths[i]; ok {
			width = w
		}

		if colOffset <= t.offset.X-width-padding {
			// before visible content
		} else if colOffset <= headWidth || colOffset <= t.offset.X {
			minCol = i
			offX = colOffset
		}
		if colOffset < t.offset.X+size.Width {
			maxCol = i + 1
		}

		colOffset += width + padding
		visible[i] = width
	}
	return
}

func (t *StaticTable) stickyRowHeights(rowHeight float32, rows int) (visible []float32) {
	if rows == 0 {
		return []float32{}
	}

	max := t.StickyRowCount
	if max > rows {
		max = rows
	}

	visible = make([]float32, max)

	if len(t.rowHeights) == 0 {
		for i := 0; i < max; i++ {
			visible[i] = rowHeight
		}
		return
	}

	for i := 0; i < max; i++ {
		height := rowHeight

		if h, ok := t.rowHeights[i]; ok {
			height = h
		}

		visible[i] = height
	}
	return
}

func (t *StaticTable) visibleRowHeights(rowHeight float32, rows int) (visible map[int]float32, offY float32, minRow, maxRow int) {
	maxRow = rows
	rowOffset, headHeight := float32(0), float32(0)
	isVisible := false
	visible = make(map[int]float32)

	if t.content.Size().Height <= 0 {
		return
	}

	padding := t.Theme().Size(theme.SizeNamePadding)
	stick := t.StickyRowCount
	size := t.size.Load()

	if len(t.rowHeights) == 0 {
		paddedHeight := rowHeight + padding

		offY = float32(math.Floor(float64(t.offset.Y/paddedHeight))) * paddedHeight
		minRow = int(math.Floor(float64(offY / paddedHeight)))
		//maxRow = int(math.Ceil(float64((t.offset.Y + size.Height) / paddedHeight)))

		if minRow > rows-1 {
			minRow = rows - 1
		}
		if minRow < 0 {
			minRow = 0
		}

		if maxRow > rows {
			maxRow = rows
		}

		visible = make(map[int]float32, maxRow-minRow+stick)
		for i := 0; i < maxRow; i++ {
			visible[i] = rowHeight
		}
		for i := 0; i < stick; i++ {
			visible[i] = rowHeight
		}
		return
	}

	for i := 0; i < rows; i++ {
		height := rowHeight
		if h, ok := t.rowHeights[i]; ok {
			height = h
		}

		if rowOffset <= t.offset.Y-height-padding {
			// before visible content
		} else if rowOffset <= headHeight || rowOffset <= t.offset.Y {
			minRow = i
			offY = rowOffset
			isVisible = true
		}
		if rowOffset < t.offset.Y+size.Height {
			maxRow = i + 1
		}

		rowOffset += height + padding
		if isVisible || i < stick {
			visible[i] = height
		}
	}
	return
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*staticTableRender)(nil)

type staticTableRender struct {
	widget.BaseRenderer
	t *StaticTable
}

func (t *staticTableRender) Layout(s fyne.Size) {
	th := t.t.Theme()
	t.t.propertyLock.RLock()

	t.calculateHeaderSizes(th)
	off := fyne.NewPos(t.t.stuckWidth, t.t.stuckHeight)
	if t.t.ShowHeaderRow {
		off.Y += t.t.headerSize.Height
	}
	if t.t.ShowHeaderColumn {
		off.X += t.t.headerSize.Width
	}
	t.t.propertyLock.RUnlock()

	t.t.content.Move(off)
	t.t.content.Resize(s.SubtractWidthHeight(off.X, off.Y))

	t.t.top.Move(fyne.NewPos(off.X, 0))
	t.t.top.Resize(fyne.NewSize(s.Width-off.X, off.Y))
	t.t.left.Move(fyne.NewPos(0, off.Y))
	t.t.left.Resize(fyne.NewSize(off.X, s.Height-off.Y))
	t.t.corner.Resize(fyne.NewSize(off.X, off.Y))
	t.t.dividerLayer.Resize(s)
	t.t.dividerLayer.Show()

}

func (t *staticTableRender) MinSize() fyne.Size {
	sep := t.t.Theme().Size(theme.SizeNamePadding)
	t.t.propertyLock.RLock()
	defer t.t.propertyLock.RUnlock()

	min := t.t.content.MinSize().Max(t.t.cellSize)
	if t.t.ShowHeaderRow {
		min.Height += t.t.headerSize.Height + sep
	}
	if t.t.ShowHeaderColumn {
		min.Width += t.t.headerSize.Width + sep
	}
	if t.t.StickyRowCount > 0 {
		for i := 0; i < t.t.StickyRowCount; i++ {
			height := t.t.cellSize.Height
			if h, ok := t.t.rowHeights[i]; ok {
				height = h
			}

			min.Height += height + sep
		}
	}
	if t.t.StickyColumnCount > 0 {
		for i := 0; i < t.t.StickyColumnCount; i++ {
			width := t.t.cellSize.Width
			if w, ok := t.t.columnWidths[i]; ok {
				width = w
			}

			min.Width += width + sep
		}
	}
	return min
}

func (t *staticTableRender) Refresh() {
	th := t.t.Theme()
	t.t.propertyLock.Lock()
	t.t.headerSize = t.t.createHeader().MinSize()
	if t.t.columnWidths != nil {
		if v, ok := t.t.columnWidths[-1]; ok {
			t.t.headerSize.Width = v
		}
	}
	if t.t.rowHeights != nil {
		if v, ok := t.t.rowHeights[-1]; ok {
			t.t.headerSize.Height = v
		}
	}
	t.t.cellSize = t.t.templateSize()
	t.calculateHeaderSizes(th)
	t.t.propertyLock.Unlock()

	t.Layout(t.t.Size())
	t.t.cells.Refresh()
}

func (t *staticTableRender) calculateHeaderSizes(th fyne.Theme) {
	t.t.stuckXOff = 0
	t.t.stuckYOff = 0

	if t.t.ShowHeaderRow {
		t.t.stuckYOff = t.t.headerSize.Height
	}
	if t.t.ShowHeaderColumn {
		t.t.stuckXOff = t.t.headerSize.Width
	}

	separatorThickness := th.Size(theme.SizeNamePadding)
	stickyColWidths := t.t.stickyColumnWidths(t.t.cellSize.Width, t.t.StickyColumnCount)
	stickyRowHeights := t.t.stickyRowHeights(t.t.cellSize.Height, t.t.StickyRowCount)

	var stuckHeight float32
	for _, rowHeight := range stickyRowHeights {
		stuckHeight += rowHeight + separatorThickness
	}
	t.t.stuckHeight = stuckHeight
	var stuckWidth float32
	for _, colWidth := range stickyColWidths {
		stuckWidth += colWidth + separatorThickness
	}
	t.t.stuckWidth = stuckWidth
}

// Declare conformity with Widget interface.
var _ fyne.Widget = (*headerTableCells)(nil)

type headerTableCells struct {
	BaseWidget
	t *StaticTable
}

func newHeaderTableCells(t *StaticTable) *headerTableCells {
	c := &headerTableCells{t: t}
	c.ExtendBaseWidget(c)
	return c
}

func (c *headerTableCells) CreateRenderer() fyne.WidgetRenderer {
	th := c.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	r := &headerTableCellsRenderer{cells: c, pool: &syncPool{}, headerPool: &syncPool{},
		visible: make(map[TableCellID]fyne.CanvasObject), headers: make(map[TableCellID]fyne.CanvasObject),
		headRowBG: canvas.NewRectangle(th.Color(theme.ColorNameHeaderBackground, v)), headColBG: canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground)),
		headRowStickyBG: canvas.NewRectangle(th.Color(theme.ColorNameHeaderBackground, v)), headColStickyBG: canvas.NewRectangle(theme.Color(theme.ColorNameHeaderBackground)),
	}

	c.t.moveCallback = r.moveIndicators
	return r
}

func (c *headerTableCells) Resize(s fyne.Size) {
	c.BaseWidget.Resize(s)
	c.Refresh() // trigger a redraw
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*headerTableCellsRenderer)(nil)

type headerTableCellsRenderer struct {
	widget.BaseRenderer

	init             bool
	cells            *headerTableCells
	pool, headerPool pool
	visible, headers map[TableCellID]fyne.CanvasObject
	dividers         []fyne.CanvasObject
	drawAllCells     bool

	headColBG, headRowBG, headRowStickyBG, headColStickyBG *canvas.Rectangle
}

func (r *headerTableCellsRenderer) Layout(fyne.Size) {
	r.cells.propertyLock.Lock()
	r.moveIndicators()
	r.cells.propertyLock.Unlock()
}

func (r *headerTableCellsRenderer) MinSize() fyne.Size {
	r.cells.propertyLock.RLock()
	defer r.cells.propertyLock.RUnlock()
	rows, cols := 0, 0
	if f := r.cells.t.Length; f != nil {
		rows, cols = r.cells.t.Length()
	} else {
		fyne.LogError("Missing Length callback required for StaticTable", nil)
	}

	stickRows := r.cells.t.StickyRowCount
	stickCols := r.cells.t.StickyColumnCount

	width := float32(0)
	if len(r.cells.t.columnWidths) == 0 {
		width = r.cells.t.cellSize.Width * float32(cols-stickCols)
	} else {
		cellWidth := r.cells.t.cellSize.Width
		for col := stickCols; col < cols; col++ {
			colWidth, ok := r.cells.t.columnWidths[col]
			if ok {
				width += colWidth
			} else {
				width += cellWidth
			}
		}
	}

	height := float32(0)
	if len(r.cells.t.rowHeights) == 0 {
		height = r.cells.t.cellSize.Height * float32(rows-stickRows)
	} else {
		cellHeight := r.cells.t.cellSize.Height
		for row := stickRows; row < rows; row++ {
			rowHeight, ok := r.cells.t.rowHeights[row]
			if ok {
				height += rowHeight
			} else {
				height += cellHeight
			}
		}
	}

	separatorSize := r.cells.t.Theme().Size(theme.SizeNamePadding)
	return fyne.NewSize(width+float32(cols-stickCols-1)*separatorSize, height+float32(rows-stickRows-1)*separatorSize)
}

func (r *headerTableCellsRenderer) Refresh() {
	if r.init == true {
		return
	}
	th := r.cells.t.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	r.cells.propertyLock.Lock()
	separatorThickness := th.Size(theme.SizeNamePadding)
	dataRows, dataCols := 0, 0
	if f := r.cells.t.Length; f != nil {
		dataRows, dataCols = r.cells.t.Length()
	}
	visibleColWidths, offX, _, _ := r.cells.t.visibleColumnWidths(r.cells.t.cellSize.Width, dataCols)
	if len(visibleColWidths) == 0 && dataCols > 0 { // we can't show anything until we have some dimensions
		r.cells.propertyLock.Unlock()
		return
	}
	visibleRowHeights, offY, _, _ := r.cells.t.visibleRowHeights(r.cells.t.cellSize.Height, dataRows)
	if len(visibleRowHeights) == 0 && dataRows > 0 { // we can't show anything until we have some dimensions
		r.cells.propertyLock.Unlock()
		return
	}

	updateCell := r.cells.t.UpdateCell
	if updateCell == nil {
		fyne.LogError("Missing UpdateCell callback required for StaticTable", nil)
	}

	var cellXOffset, cellYOffset float32
	startRow := 0
	startCol := 0

	wasVisible := r.visible
	r.visible = make(map[TableCellID]fyne.CanvasObject)
	var cells []fyne.CanvasObject
	displayCol := func(row, col int, rowHeight float32, cells *[]fyne.CanvasObject) {
		id := TableCellID{row, col}
		colWidth := visibleColWidths[col]
		c, ok := wasVisible[id]
		if !ok {
			c = r.pool.Obtain()
			if f := r.cells.t.CreateCell; f != nil && c == nil {
				c = createItemAndApplyThemeScope(f, r.cells.t)
			}
			if c == nil {
				return
			}
		}

		c.Move(fyne.NewPos(cellXOffset, cellYOffset))
		c.Resize(fyne.NewSize(colWidth, rowHeight))

		r.visible[id] = c
		*cells = append(*cells, c)
		cellXOffset += colWidth + separatorThickness
	}

	displayRow := func(row int, cells *[]fyne.CanvasObject) {
		rowHeight := visibleRowHeights[row]
		cellXOffset = 0

		for col := startCol; col < dataCols; col++ {
			displayCol(row, col, rowHeight, cells)
		}
		cellXOffset = r.cells.t.content.Offset.X
		stick := r.cells.t.StickyColumnCount
		if r.cells.t.ShowHeaderColumn {
			cellXOffset += r.cells.t.headerSize.Width
			stick--
		}
		cellYOffset += rowHeight + separatorThickness
	}

	cellYOffset = 0
	for row := startRow; row < dataRows; row++ {
		displayRow(row, &cells)
	}

	inline := r.refreshHeaders(visibleRowHeights, visibleColWidths, offX, offY, startRow, dataRows, startCol, dataCols,
		separatorThickness, th, v)
	cells = append(cells, inline...)

	offX -= r.cells.t.content.Offset.X
	cellYOffset = r.cells.t.stuckYOff

	for id, old := range wasVisible {
		if _, ok := r.visible[id]; !ok {
			r.pool.Release(old)
		}
	}
	visible := r.visible
	headers := r.headers

	r.cells.propertyLock.Unlock()
	r.SetObjects(cells)

	if updateCell != nil {
		for id, cell := range visible {

			updateCell(id, cell)
		}
	}
	for id, head := range headers {
		r.cells.t.updateHeader(id, head)
	}

	r.moveIndicators()
	r.init = true
}
func (r *headerTableCellsRenderer) scroll() {
	th := r.cells.t.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	r.cells.propertyLock.Lock()
	separatorThickness := th.Size(theme.SizeNamePadding)
	dataRows, dataCols := 0, 0
	if f := r.cells.t.Length; f != nil {
		dataRows, dataCols = r.cells.t.Length()
	}

	visibleColWidths, offX, minCol, _ := r.cells.t.visibleColumnWidths(r.cells.t.cellSize.Width, dataCols)
	if len(visibleColWidths) == 0 && dataCols > 0 { // we can't show anything until we have some dimensions
		r.cells.propertyLock.Unlock()
		return
	}
	visibleRowHeights, offY, minRow, _ := r.cells.t.visibleRowHeights(r.cells.t.cellSize.Height, dataRows)
	if len(visibleRowHeights) == 0 && dataRows > 0 { // we can't show anything until we have some dimensions
		r.cells.propertyLock.Unlock()
		return
	}
	updateCell := r.cells.t.UpdateCell
	if updateCell == nil {
		fyne.LogError("Missing UpdateCell callback required for StaticTable", nil)
	}
	_ = r.refreshHeaders(visibleRowHeights, visibleColWidths, offX, offY, minRow, dataRows, minCol, dataCols,
		separatorThickness, th, v)
	r.cells.propertyLock.Unlock()
	headers := r.headers
	for id, head := range headers {
		r.cells.t.updateHeader(id, head)
	}
	r.moveIndicators()
}

func (r *headerTableCellsRenderer) refreshForID(toDraw TableCellID) {
	th := r.cells.t.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()

	r.cells.propertyLock.Lock()
	separatorThickness := th.Size(theme.SizeNamePadding)
	dataRows, dataCols := 0, 0
	if f := r.cells.t.Length; f != nil {
		dataRows, dataCols = r.cells.t.Length()
	}
	visibleColWidths, offX, minCol, maxCol := r.cells.t.visibleColumnWidths(r.cells.t.cellSize.Width, dataCols)
	if len(visibleColWidths) == 0 && dataCols > 0 { // we can't show anything until we have some dimensions
		r.cells.propertyLock.Unlock()
		return
	}
	visibleRowHeights, offY, minRow, maxRow := r.cells.t.visibleRowHeights(r.cells.t.cellSize.Height, dataRows)
	if len(visibleRowHeights) == 0 && dataRows > 0 { // we can't show anything until we have some dimensions
		r.cells.propertyLock.Unlock()
		return
	}

	updateCell := r.cells.t.UpdateCell
	if updateCell == nil {
		fyne.LogError("Missing UpdateCell callback required for StaticTable", nil)
	}

	var cellXOffset, cellYOffset float32
	stickRows := r.cells.t.StickyRowCount
	if r.cells.t.ShowHeaderRow {
		cellYOffset += r.cells.t.headerSize.Height
	}
	stickCols := r.cells.t.StickyColumnCount
	if r.cells.t.ShowHeaderColumn {
		cellXOffset += r.cells.t.headerSize.Width
	}
	startRow := minRow + stickRows
	if startRow < stickRows {
		startRow = stickRows
	}
	startCol := minCol + stickCols
	if startCol < stickCols {
		startCol = stickCols
	}

	wasVisible := r.visible
	r.visible = make(map[TableCellID]fyne.CanvasObject)
	var cells []fyne.CanvasObject
	displayCol := func(row, col int, rowHeight float32, cells *[]fyne.CanvasObject) {
		id := TableCellID{row, col}
		colWidth := visibleColWidths[col]
		c, ok := wasVisible[id]
		if !ok {
			c = r.pool.Obtain()
			if f := r.cells.t.CreateCell; f != nil && c == nil {
				c = createItemAndApplyThemeScope(f, r.cells.t)
			}
			if c == nil {
				return
			}
		}

		c.Move(fyne.NewPos(cellXOffset, cellYOffset))
		c.Resize(fyne.NewSize(colWidth, rowHeight))

		r.visible[id] = c
		*cells = append(*cells, c)
		cellXOffset += colWidth + separatorThickness
	}

	displayRow := func(row int, cells *[]fyne.CanvasObject) {
		rowHeight := visibleRowHeights[row]
		cellXOffset = offX

		for col := startCol; col < maxCol; col++ {
			displayCol(row, col, rowHeight, cells)
		}
		cellXOffset = r.cells.t.content.Offset.X
		stick := r.cells.t.StickyColumnCount
		if r.cells.t.ShowHeaderColumn {
			cellXOffset += r.cells.t.headerSize.Width
			stick--
		}
		cellYOffset += rowHeight + separatorThickness
	}

	cellYOffset = offY
	for row := startRow; row < maxRow; row++ {
		displayRow(row, &cells)
	}

	inline := r.refreshHeaders(visibleRowHeights, visibleColWidths, offX, offY, startRow, maxRow, startCol, maxCol,
		separatorThickness, th, v)
	cells = append(cells, inline...)

	offX -= r.cells.t.content.Offset.X
	cellYOffset = r.cells.t.stuckYOff
	for row := 0; row < stickRows; row++ {
		displayRow(row, &r.cells.t.top.Content.(*fyne.Container).Objects)
	}

	cellYOffset = offY - r.cells.t.content.Offset.Y
	for row := startRow; row < maxRow; row++ {
		cellXOffset = r.cells.t.stuckXOff
		rowHeight := visibleRowHeights[row]
		for col := 0; col < stickCols; col++ {
			displayCol(row, col, rowHeight, &r.cells.t.left.Content.(*fyne.Container).Objects)
		}
		cellYOffset += rowHeight + separatorThickness
	}

	cellYOffset = r.cells.t.stuckYOff
	for row := 0; row < stickRows; row++ {
		cellXOffset = r.cells.t.stuckXOff
		rowHeight := visibleRowHeights[row]
		for col := 0; col < stickCols; col++ {
			displayCol(row, col, rowHeight, &r.cells.t.corner.Content.(*fyne.Container).Objects)
		}
		cellYOffset += rowHeight + separatorThickness
	}

	for id, old := range wasVisible {
		if _, ok := r.visible[id]; !ok {
			r.pool.Release(old)
		}
	}
	visible := r.visible

	r.cells.propertyLock.Unlock()
	r.SetObjects(cells)

	if updateCell != nil {
		for id, cell := range visible {
			if toDraw != allTableCellsID && toDraw != id {
				continue
			}

			updateCell(id, cell)
		}
	}
	r.moveIndicators()
}

func (r *headerTableCellsRenderer) moveIndicators() {
	rows, cols := 0, 0
	if f := r.cells.t.Length; f != nil {
		rows, cols = r.cells.t.Length()
	}
	visibleColWidths, offX, minCol, maxCol := r.cells.t.visibleColumnWidths(r.cells.t.cellSize.Width, cols)
	visibleRowHeights, offY, minRow, maxRow := r.cells.t.visibleRowHeights(r.cells.t.cellSize.Height, rows)
	th := r.cells.t.Theme()
	separatorThickness := th.Size(theme.SizeNameSeparatorThickness)
	padding := th.Size(theme.SizeNamePadding)
	dividerOff := (padding - separatorThickness) / 2

	if r.cells.t.ShowHeaderColumn {
		offX += r.cells.t.headerSize.Width
	}
	if r.cells.t.ShowHeaderRow {
		offY += r.cells.t.headerSize.Height
	}

	colDivs := maxCol - minCol - 1
	if colDivs < 0 {
		colDivs = 0
	}
	rowDivs := maxRow - minRow - 1
	if rowDivs < 0 {
		rowDivs = 0
	}

	if colDivs < 0 {
		colDivs = 0
	}
	if rowDivs < 0 {
		rowDivs = 0
	}

	if len(r.dividers) < colDivs+rowDivs {
		for i := len(r.dividers); i < colDivs+rowDivs; i++ {
			r.dividers = append(r.dividers, NewSeparator())
		}

		var objs []fyne.CanvasObject
		r.cells.t.dividerLayer.Content.(*fyne.Container).Objects = append(objs, r.dividers...)
		r.cells.t.dividerLayer.Content.Refresh()
	}

	size := r.cells.t.size.Load()

	divs := 0
	i := 0
	i = minCol
	for x := offX + r.cells.t.stuckWidth + visibleColWidths[i]; i < maxCol-1 && divs < colDivs; x += visibleColWidths[i] + padding {
		i++

		xPos := x - r.cells.t.content.Offset.X + dividerOff
		r.dividers[divs].Resize(fyne.NewSize(separatorThickness, size.Height))
		r.dividers[divs].Move(fyne.NewPos(xPos, 0))
		r.dividers[divs].Show()
		divs++
	}

	i = 0
	i = minRow
	for y := offY + r.cells.t.stuckHeight + visibleRowHeights[i]; i < maxRow-1 && divs-colDivs < rowDivs; y += visibleRowHeights[i] + padding {
		i++

		yPos := y - r.cells.t.content.Offset.Y + dividerOff
		r.dividers[divs].Resize(fyne.NewSize(size.Width, separatorThickness))
		r.dividers[divs].Move(fyne.NewPos(0, yPos))
		r.dividers[divs].Show()
		divs++
	}

	for i := divs; i < len(r.dividers); i++ {
		r.dividers[i].Hide()
	}
}

func (r *headerTableCellsRenderer) refreshHeaders(visibleRowHeights, visibleColWidths map[int]float32, offX, offY float32,
	startRow, maxRow, startCol, maxCol int, separatorThickness float32, th fyne.Theme, v fyne.ThemeVariant) []fyne.CanvasObject {
	wasVisible := r.headers
	r.headers = make(map[TableCellID]fyne.CanvasObject)
	headerMin := r.cells.t.headerSize
	rowHeight := headerMin.Height
	colWidth := headerMin.Width

	var cells, over []fyne.CanvasObject
	over = []fyne.CanvasObject{r.headRowBG}
	if r.cells.t.ShowHeaderRow {
		cellXOffset := offX - r.cells.t.content.Offset.X
		displayColHeader := func(col int, list *[]fyne.CanvasObject) {
			id := TableCellID{-1, col}
			colWidth := visibleColWidths[col]
			c, ok := wasVisible[id]
			if !ok {
				c = r.headerPool.Obtain()
				if c == nil {
					c = r.cells.t.createHeader()
				}
				if c == nil {
					return
				}
			}

			c.Move(fyne.NewPos(cellXOffset, 0))
			c.Resize(fyne.NewSize(colWidth, rowHeight))

			r.headers[id] = c
			*list = append(*list, c)
			cellXOffset += colWidth + separatorThickness
		}
		for col := startCol; col < maxCol; col++ {
			displayColHeader(col, &over)
		}
	}
	r.cells.t.top.Content.(*fyne.Container).Objects = over
	r.cells.t.top.Content.Refresh()

	over = []fyne.CanvasObject{r.headColBG}
	if r.cells.t.ShowHeaderColumn {
		cellYOffset := offY - r.cells.t.content.Offset.Y
		displayRowHeader := func(row int, list *[]fyne.CanvasObject) {
			id := TableCellID{row, -1}
			rowHeight := visibleRowHeights[row]
			c, ok := wasVisible[id]
			if !ok {
				c = r.headerPool.Obtain()
				if c == nil {
					c = r.cells.t.createHeader()
				}
				if c == nil {
					return
				}
			}

			c.Move(fyne.NewPos(0, cellYOffset))
			c.Resize(fyne.NewSize(colWidth, rowHeight))

			r.headers[id] = c
			*list = append(*list, c)
			cellYOffset += rowHeight + separatorThickness
		}
		for row := startRow; row < maxRow; row++ {
			displayRowHeader(row, &over)
		}
	}
	r.cells.t.left.Content.(*fyne.Container).Objects = over
	r.cells.t.left.Content.Refresh()
	for id, old := range wasVisible {
		if _, ok := r.headers[id]; !ok {
			r.headerPool.Release(old)
		}
	}
	return cells
}

type staticClip struct {
	widget.Scroll

	t *StaticTable
}

func newHeaderClip(t *StaticTable, o fyne.CanvasObject) *staticClip {
	c := &staticClip{t: t}
	c.Content = o
	c.Direction = widget.ScrollNone

	return c
}
