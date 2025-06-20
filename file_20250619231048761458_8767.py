import collections

CARD_RANKS = {
    '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'T': 10,
    'J': 11, 'Q': 12, 'K': 13, 'A': 14
}

def parse_card(card_str):
    rank_char = card_str[0]
    suit_char = card_str[1]
    return CARD_RANKS[rank_char], suit_char

def evaluate_hand(hand):
    parsed_cards = [parse_card(card) for card in hand]
    
    ranks = sorted([card[0] for card in parsed_cards])
    suits = [card[1] for card in parsed_cards]
    
    rank_counts = collections.Counter(ranks)
    num_unique_ranks = len(rank_counts)
    
    is_flush = len(set(suits)) == 1
    
    is_straight = False
    original_ranks = list(ranks) # Keep original for non-Ace-low straight check
    if num_unique_ranks == 5:
        # Check for normal straight (consecutive ranks)
        if ranks[4] - ranks[0] == 4:
            is_straight = True
        # Check for Ace-low straight (A,2,3,4,5)
        elif ranks == [2, 3, 4, 5, 14]: # 14 is Ace
            is_straight = True
            ranks = [1, 2, 3, 4, 5] # Re-map Ace to 1 for comparison purposes
    
    # Hand types in descending order of value (higher tuple is better hand)
    
    # Straight Flush (including Royal Flush)
    if is_straight and is_flush:
        if ranks == [1,2,3,4,5]: # Ace-low straight flush
            return (9, 5) # Use 5 as the high card for A-5 straight
        return (9, ranks[4]) # (Hand_Type_Score, High_Card_Rank)
    
    # Four of a Kind
    if 4 in rank_counts.values():
        four_kind_rank = 0
        kicker_rank = 0
        for rank, count in rank_counts.items():
            if count == 4:
                four_kind_rank = rank
            else:
                kicker_rank = rank
        return (8, four_kind_rank, kicker_rank)
    
    # Full House
    if 3 in rank_counts.values() and 2 in rank_counts.values():
        three_kind_rank = 0
        pair_rank = 0
        for rank, count in rank_counts.items():
            if count == 3:
                three_kind_rank = rank
            elif count == 2:
                pair_rank = rank
        return (7, three_kind_rank, pair_rank)
    
    # Flush
    if is_flush:
        return (6, *sorted(original_ranks, reverse=True)) # (Hand_Type_Score, Ranks_Descending)
    
    # Straight
    if is_straight:
        if ranks == [1,2,3,4,5]: # Ace-low straight
            return (5, 5) # Use 5 as the high card for A-5 straight
        return (5, ranks[4]) # (Hand_Type_Score, High_Card_Rank)
    
    # Three of a Kind
    if 3 in rank_counts.values():
        three_kind_rank = 0
        kickers = []
        for rank, count in rank_counts.items():
            if count == 3:
                three_kind_rank = rank
            else:
                kickers.append(rank)
        return (4, three_kind_rank, *sorted(kickers, reverse=True))
    
    # Two Pair
    if list(rank_counts.values()).count(2) == 2:
        pairs = []
        kicker = 0
        for rank, count in rank_counts.items():
            if count == 2:
                pairs.append(rank)
            else:
                kicker = rank
        pairs.sort(reverse=True)
        return (3, pairs[0], pairs[1], kicker)
    
    # One Pair
    if 2 in rank_counts.values():
        pair_rank = 0
        kickers = []
        for rank, count in rank_counts.items():
            if count == 2:
                pair_rank = rank
            else:
                kickers.append(rank)
        return (2, pair_rank, *sorted(kickers, reverse=True))
    
    # High Card
    return (1, *sorted(original_ranks, reverse=True))

# Additional implementation at 2025-06-19 23:12:09
import collections

def parse_card(card_str):
    rank_str = card_str[:-1]
    suit = card_str[-1]
    if rank_str == 'T':
        rank = 10
    elif rank_str == 'J':
        rank = 11
    elif rank_str == 'Q':
        rank = 12
    elif rank_str == 'K':
        rank = 13
    elif rank_str == 'A':
        rank = 14
    else:
        rank = int(rank_str)
    return (rank, suit)

def parse_hand(hand_strs):
    return [parse_card(s) for s in hand_strs]

def get_ranks(hand):
    return sorted([card[0] for card in hand], reverse=True)

def get_suits(hand):
    return [card[1] for card in hand]

def is_flush(hand):
    suits = get_suits(hand)
    return len(set(suits)) == 1

def is_straight(hand):
    ranks = get_ranks(hand)
    if len(set(ranks)) == 5 and ranks[0] - ranks[4] == 4:
        return True
    if set(ranks) == {14, 5, 4, 3, 2}:
        return True
    return False

def get_rank_counts(hand):
    ranks = get_ranks(hand)
    return collections.Counter(ranks)

def evaluate_hand(hand_strs):
    hand = parse_hand(hand_strs)
    ranks = get_ranks(hand)
    rank_counts = get_rank_counts(hand)

    is_hand_flush = is_flush(hand)
    is_hand_straight = is_straight(hand)

    counts = sorted(rank_counts.values(), reverse=True)
    unique_ranks = sorted(rank_counts.keys(), reverse=True)

    if is_hand_straight and is_hand_flush:
        if set(ranks) == {14, 5, 4, 3, 2}:
            return (8, 5)
        return (8, ranks[0])

    if counts[0] == 4:
        quad_rank = unique_ranks[0] if rank_counts[unique_ranks[0]] == 4 else unique_ranks[1]
        kicker = [r for r in unique_ranks if r != quad_rank][0]
        return (7, quad_rank, kicker)

    if counts[0] == 3 and counts[1] == 2:
        trip_rank = unique_ranks[0] if rank_counts[unique_ranks[0]] == 3 else unique_ranks[1]
        pair_rank = unique_ranks[0] if rank_counts[unique_ranks[0]] == 2 else unique_ranks[1]
        return (6, trip_rank, pair_rank)

    if is_hand_flush:
        return (5, *ranks)

    if is_hand_straight:
        if set(ranks) == {14, 5, 4, 3, 2}:
            return (4, 5)
        return (4, ranks[0])

    if counts[0] == 3:
        trip_rank = unique_ranks[0] if rank_counts[unique_ranks[0]] == 3 else unique_ranks[1]
        kickers = sorted([r for r in unique_ranks if r != trip_rank], reverse=True)
        return (3, trip_rank, *kickers)

    if counts[0] == 2 and counts[1] == 2:
        pair_ranks = sorted([r for r, c in rank_counts.items() if c == 2], reverse=True)
        kicker = [r for r, c in rank_counts.items() if c == 1][0]
        return (2, pair_ranks[0], pair_ranks[1], kicker)

    if counts[0] == 2:
        pair_rank = [r for r, c in rank_counts.items() if c == 2][0]
        kickers = sorted([r for r, c in rank_counts.items() if c == 1], reverse=True)
        return (1, pair_rank, *kickers)

    return (0, *ranks)