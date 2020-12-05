/**
 * sourcemod/menus.go
 * 
 * Copyright 2020 Nirari Technologies, Alliedmodders LLC.
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 * 
 */

package main


type MenuStyle int
const (
	MenuStyle_Default = MenuStyle(0)      /**< The "default" menu style for the mod */
	MenuStyle_Valve        /**< The Valve provided menu style (Used on HL2DM) */
	MenuStyle_Radio         /**< The simpler menu style commonly used on CS:S */
)


type MenuAction int
const (
	MenuAction_Start = (1<<0)      /**< A menu has been started (nothing passed) */
	MenuAction_Display = (1<<1)    /**< A menu is about to be displayed (param1=client, param2=MenuPanel Handle) */
	MenuAction_Select = (1<<2)     /**< An item was selected (param1=client, param2=item) */
	MenuAction_Cancel = (1<<3)     /**< The menu was cancelled (param1=client, param2=reason) */
	MenuAction_End = (1<<4)        /**< A menu display has fully ended.
                                         param1 is the MenuEnd reason, and if it's MenuEnd_Cancelled, then
                                         param2 is the MenuCancel reason from MenuAction_Cancel. */
	MenuAction_VoteEnd = (1<<5)    /**< (VOTE ONLY): A vote sequence has succeeded (param1=chosen item)
                                         This is not called if SetVoteResultCallback has been used on the menu. */
	MenuAction_VoteStart = (1<<6)  /**< (VOTE ONLY): A vote sequence has started (nothing passed) */
	MenuAction_VoteCancel = (1<<7) /**< (VOTE ONLY): A vote sequence has been cancelled (param1=reason) */
	MenuAction_DrawItem = (1<<8)   /**< An item is being drawn; return the new style (param1=client, param2=item) */
	MenuAction_DisplayItem = (1<<9) /**< Item text is being drawn to the display (param1=client, param2=item)
                                         To change the text, use RedrawMenuItem().
                                         If you do so, return its return value.  Otherwise, return 0. */
)


const (
	MENU_ACTIONS_DEFAULT =    MenuAction_Select|MenuAction_Cancel|MenuAction_End
/** All menu actions */
	MENU_ACTIONS_ALL =        MenuAction(0xFFFFFFFF)

	MENU_NO_PAGINATION =      0           /**< Menu should not be paginated (10 items max) */
	MENU_TIME_FOREVER =       0           /**< Menu should be displayed as long as possible */

	ITEMDRAW_DEFAULT =            (0)     /**< Item should be drawn normally */
	ITEMDRAW_DISABLED =           (1<<0)  /**< Item is drawn but not selectable */
	ITEMDRAW_RAWLINE =            (1<<1)  /**< Item should be a raw line, without a slot */
	ITEMDRAW_NOTEXT =             (1<<2)  /**< No text should be drawn */
	ITEMDRAW_SPACER =             (1<<3)  /**< Item should be drawn as a spacer, if possible */
	ITEMDRAW_IGNORE =     ((1<<1)|(1<<2)) /**< Item should be completely ignored (rawline + notext) */
	ITEMDRAW_CONTROL =            (1<<4)  /**< Item is control text (back/next/exit) */

	MENUFLAG_BUTTON_EXIT =        (1<<0)  /**< Menu has an "exit" button (default if paginated) */
	MENUFLAG_BUTTON_EXITBACK =    (1<<1)  /**< Menu has an "exit back" button */
	MENUFLAG_NO_SOUND =           (1<<2)  /**< Menu will not have any select sounds */
	MENUFLAG_BUTTON_NOVOTE =      (1<<3)  /**< Menu has a "No Vote" button at slot 1 */

	VOTEINFO_CLIENT_INDEX =       0       /**< Client index */
	VOTEINFO_CLIENT_ITEM =        1       /**< Item the client selected, or -1 for none */
	VOTEINFO_ITEM_INDEX =         0       /**< Item index */
	VOTEINFO_ITEM_VOTES =         1       /**< Number of votes for the item */

	VOTEFLAG_NO_REVOTES =         (1<<0)  /**< Players cannot change their votes */

	MenuCancel_Disconnected = -1   /**< Client dropped from the server */
	MenuCancel_Interrupted = -2    /**< Client was interrupted with another menu */
	MenuCancel_Exit = -3           /**< Client exited via "exit" */
	MenuCancel_NoDisplay = -4      /**< Menu could not be displayed to the client */
	MenuCancel_Timeout = -5        /**< Menu timed out */
	MenuCancel_ExitBack = -6       /**< Client selected "exit back" on a paginated menu */
	
	VoteCancel_Generic = -1
	VoteCancel_NoVotes = -2
	
	MenuEnd_Selected = 0           /**< Menu item was selected */
	MenuEnd_VotingDone = -1        /**< Voting finished */
	MenuEnd_VotingCancelled = -2   /**< Voting was cancelled */
	MenuEnd_Cancelled = -3         /**< Menu was cancelled (reason in param2) */
	MenuEnd_Exit = -4              /**< Menu was cleanly exited via "exit" */
	MenuEnd_ExitBack = -5          /**< Menu was cleanly exited via "back" */
)


type MenuSource int
const (
	MenuSource_None = MenuSource(0)            /**< No menu is being displayed */
	MenuSource_External        /**< External menu */
	MenuSource_Normal          /**< A basic menu is being displayed */
	MenuSource_RawPanel         /**< A display is active, but it is not tied to a menu */
)


type (
	MenuHandler func(menu Menu, act MenuAction, parm1, parm2 int) int
	/// new Panel(Handle hStyle = null);
	Panel struct {
		TextRemaining, CurrentKey int
		Style Handle
	}
	/// new Menu(MenuHandler handler, MenuAction actions=MENU_ACTIONS_DEFAULT);
	Menu struct {
		Pagination, OptionFlags, ItemCount, Selection int
		ExitButton, ExitBackButton, NoVoteButton bool
		Style Handle
		VoteResultCallback VoteHandler
	}
)

func (Panel) SetTitle(text string, onlyIfEmpty bool)
func (Panel) DrawItem(text string, style int) int
func (Panel) DrawText(text string) bool
func (Panel) CanDrawFlags(style int) bool
func (Panel) SetKeys(keys int) bool
func (Panel) Send(client int, handler MenuHandler, time int) bool


func (Menu) Display(client, time int) bool
func (Menu) DisplayAt(client, first_item, time int) bool
func (Menu) AddItem(info, display string, style int) bool
func (Menu) InsertItem(position int, info, display string, style int) bool
func (Menu) RemoveItem(position int) bool
func (Menu) RemoveAllItems()
func (Menu) GetItem(position int, infoBuf []char, infoBufLen int, style *int, dispBuf []char, dispBufLen int) bool
func (Menu) SetTitle(fmt string, args ...any)
func (Menu) GetTitle(buffer []char, maxlength int)
func (Menu) ToPanel() Panel
func (Menu) Cancel()
func (Menu) DisplayVote(clients []int, numClients, time, flags int) bool
func (Menu) DisplayVoteToAll(time, flags int) bool


func IsVoteInProgress(menu Handle) bool
func CancelVote()


type VoteHandler func(menu Menu, num_votes, num_clients int, client_info [][]int, num_items int, item_info [][]int)

func CheckVoteDelay() int
func IsClientInVotePool(client Entity) bool
func RedrawClientVoteMenu(client Entity, revotes bool) bool
func GetClientMenu(client Entity, hStyle Handle) MenuSource
func CancelClientMenu(client Entity, autoIgnore bool, hStyle Handle) bool
func InternalShowMenu(client Entity, str string, time, keys int, handler MenuHandler) bool
func GetMenuVoteInfo(param2 int, winningVotes, totalVotes *int)
func IsNewVoteAllowed() bool