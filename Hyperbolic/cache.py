# ---------------------------------------------------------------------
# cache.py
# Author: Reuben Agogoe, Stephen Dong, and Jimmy Hoang
# ---------------------------------------------------------------------
from datetime import datetime, date
import collections

class item(): 
    def __init__(self, value):

        # currently just store a string for now 
        item.value = value

        #datetime object 
        item.ins_time = None
        item.accessed = 0 # idk lol



class Hyperbolic():

    def __init__(self, max_capacity):

        self._max_capacity = max_capacity
        
        self.current_count = 0

        # store item types 
        self.storage = [None] * max_capacity
        
        self.dictionary = {} 
    
    # private clases argmin to get the minimum item
    # input a list, defaults to zero in the case that there is no minimum 
    def argmin(list):

        minindex = 0

        for i, value in enumerate(list):
            list[i] < list[0]
            minindex = i 

        return 


    def evict_which(self):
        sampled_items = random_sample(S) # someone please implement the random_sample
        return argmin(p(i) for i in sampled_items) #someone please implement the argmin 
        # seems like it returns the minimum priority one, thisi s an area for confusion. 

    
    def evict(self):
        self.evict_which() 

    

    def when_get(item):
        item.accessed += 1

    def get(item): 
        print("test")
       

    # currently operating under the assumpption that this is the putting the item in the sampler
    def add_to_sampler(self):
        # now here we have to figure out how we want to add it to the list. 
        print("test")


    def when_put(self, item):
        item.accessed = 1
        item.ins_time = datetime.combine(date.today(), date.now().time())

        add_to_sampler(item) 
    
    def put(self, item): 

        # check for max capacity, evict value if there is too much value 
        if (self.current_count == self.max_capacity): 
            
            
            self.evict()

        # update teh necessary metadata
        self.when_put(item) 



    


    # An item's priority is an estimate of its frequency since it entered
    # the cache: p_i = n_i / t_i, where n_i is the request count
    # for i since it entered the cache, and t_i is the time since
    # it entered the cache.
    def p(item):
        
        now = datetime.now().time()

        time_in_cache = datetime.combine(date.today(), now) - item.ins_time

        return item.accessed / (time_in_cache)
        # we are supposed to evict what is hyperbolic lol 


# modify this to fit needs of hyperbolic

class LFUCache:
   def __init__(self, capacity):
      self.remain = capacity
      self.least_freq = 1

      # freq -> key -> (value, freq)
      # The main difference between defaultdict and dict is that 
      # when you try to access or modify a key that's not present 
      # in the dictionary, a default value is automatically given to that key. 
      self.node_for_freq = collections.defaultdict(collections.OrderedDict)

      # so -> orderedDict -> preserves the order inwhich tehy are inserted, and that is stored within a defaultdict 

      # key -> (value, freq)
      self.node_for_key = dict()

   def _update(self, key, value):
    # get freq of key
      _, freq = self.node_for_key[key]

      #
      self.node_for_freq[freq].pop(key)
      if len(self.node_for_freq[self.least_freq]) == 0:
         self.least_freq += 1
      self.node_for_freq[freq+1][key] = (value, freq+1)
      self.node_for_key[key] = (value, freq+1)

   def get(self, key):

    # key not found
      if key not in self.node_for_key:
         return -1

    # key found, so update dicts and return value
      value = self.node_for_key[key][0]
      self._update(key, value)
      return value

   def put(self, key, value):

    # check that there is room in the cache
    if self.remain == 0:
        # last = false means that FIFO is enforced instead of LIFO
        removed = self.node_for_freq[self.least_freq].popitem(last = False)
        self.node_for_key.pop(removed[0])

    # key exists, so update its value
    if key in self.node_for_key:
        self._update(key, value)

    # key does not exist, so add it
    else:

        # update dicts
         self.node_for_key[key] = (value, 1)
         self.node_for_freq[1][key] = (value, 1)


         if self.remain == 0:
            removed = self.node_for_freq[self.least_freq].popitem(last=False)
            self.node_for_key.pop(removed[0])
         else:
            self.remain -= 1
            self.least_freq = 1